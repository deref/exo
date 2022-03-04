package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/mathutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobWatchCmd)
}

var jobWatchCmd = &cobra.Command{
	Use:   "watch <job-id>",
	Short: "Watch a job's progress",
	Long: `Tails the events from all tasks in a job and continuously renders a task
tree until the job has finished running.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return watchJob(ctx, args[0])
	},
}

type jobEventFragment struct {
	Type       string
	Message    string
	Timestamp  scalars.Instant
	SourceType string
	SourceID   string
	Job        *jobEventJobFragment
	JobID      *string
	Task       *struct {
		ID    string
		Error *string
	}
}

type jobEventJobFragment struct {
	URL      string
	Tasks    []taskFragment
	RootTask taskFragment
}

func watchJob(ctx context.Context, jobID string) error {
	beginExclusive()
	defer endExclusive()

	type watchJobSubscription struct {
		Event jobEventFragment `graphql:"watchJob(id: $id, debug: $debug)"`
	}
	var res watchJobSubscription
	sub := api.Subscribe(ctx, svc, &res, map[string]interface{}{
		"id":    jobID,
		"debug": isDebugMode(),
	})
	defer sub.Stop()

	{
		stopSignals := make(chan os.Signal)
		signal.Notify(stopSignals, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stopSignals
			sub.Stop()
		}()
	}

	w := &EventWriter{
		W: os.Stdout,
	}
	w.Init()

	interactive := isInteractive()
	verbose := !interactive || isDebugMode()

	var panel *term.BottomPanel
	var jp *jobPrinter
	if interactive {
		panel = &term.BottomPanel{}
		defer func() {
			content := panel.Content()
			panel.Close()
			cmdutil.Show(content)
		}()

		jp = &jobPrinter{
			Spinner:            []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
			CollapseSuccessful: !verbose,
		}
	}

	var job *jobEventJobFragment

	taskProgress := make(map[string]struct {
		Reported float64
		Current  float64
	})

	// Periodically tick to keep spinner animation lively, even when
	// there are no events. Start ticking only after the first event.
	initialized := false
	var tickC <-chan time.Time

watching:
	for {
		var event jobEventFragment
		select {
		case eventInterface, ok := <-sub.Events():
			if !ok {
				break watching
			}
			event = api.OperationData(eventInterface).(jobEventFragment)
		case <-tickC:
			event = jobEventFragment{
				Type: "Tick",
			}
		}
		if event.Job != nil {
			job = event.Job
		}
		task := event.Task

		sourceID := event.SourceID
		sourceLabel := fmt.Sprintf("%s:%s", event.SourceType, sourceID)
		switch event.SourceType {
		case "System":
			sourceLabel = "system"
		case "Job":
			// We're only tracking one job, so no need to repeatedly show its name.
			sourceLabel = "job"
		case "Task":
			// For any given job, a short prefix is extremely likely to be unique.
			sourceLabel = fmt.Sprintf("task:%s", sourceID[:5])
		}

		switch event.Type {
		case "Tick":
			jp.Iteration++

		case "JobWatched":
			fmt.Println("Job URL:", job.URL)
			if initialized {
				return errors.New("already received JobWatched event")
			}
			if interactive {
				ticker := time.NewTicker(time.Second / time.Duration(len(jp.Spinner)))
				tickC = ticker.C
				defer ticker.Stop()
			}
			initialized = true

		case "JobUpdated":
			if verbose {
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, "job updated")
				// Ocassionally print task progress updates.
				const maxReportsPerTask = 20
				const reportThreshold = 100.0 / float64(maxReportsPerTask)
				for _, task := range job.Tasks {
					if task.Progress == nil {
						continue
					}
					progress := taskProgress[task.ID]
					progress.Current = task.Progress.Percent
					if progress.Current < 100 && progress.Reported+reportThreshold < progress.Current {
						progress.Reported = progress.Current
						message := fmt.Sprintf("task %s progress: %2d%%", task.ID, int(progress.Current))
						w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
					}
				}
			}

		case "TaskStarted":
			if verbose {
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, "Task Started")
			}

		case "TaskFinished":
			if verbose {
				var message string
				if task.Error == nil {
					message = "task finished; awaiting children for completion"
				} else {
					message = fmt.Sprintf("task failed: %s", *task.Error)
				}
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
			}

		case "TaskCompleted":
			if verbose {
				var message string
				if task.Error == nil {
					message = "task completed successfully"
				} else {
					message = fmt.Sprintf("task completed with failure: %s", *task.Error)
				}
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
			}

		case "JobCompleted":
			if verbose {
				message := fmt.Sprintf("job %s completed", jobID)
				w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
			}
			sub.Stop()

		case "Message":
			w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, event.Message)

		default:
			message := event.Type
			if event.Message != "" {
				message = fmt.Sprintf("%s: %s", message, event.Message)
			}
			w.PrintEvent(sourceID, event.Timestamp.GoTime(), sourceLabel, message)
		}

		if interactive {
			content := jobTreeString(jp, job.Tasks)
			height := strings.Count(content, "\n")
			panel.SetHeight(mathutil.IntMax(height, panel.Height()))
			panel.SetContent(content)
		}
	}
	if sub.Err() != nil {
		return fmt.Errorf("subscription error: %w", sub.Err())
	}

	if job == nil {
		return errors.New("never received job details")
	}
	root := job.RootTask
	if root.Error != nil {
		return fmt.Errorf("job failure: %s", *root.Error)
	}
	return nil
}
