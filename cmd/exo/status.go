package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/deref/exo/core/api"
	"github.com/deref/exo/providers/unix/components/process"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/osutil"
	psprocess "github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

type ProcessDetails struct {
	Command    string `json:"command"`
	CreateTime time.Time
}

func getProcessDetails(pid int32) ProcessDetails {

	proc, err := psprocess.NewProcess(int32(pid))
	if err != nil {
		panic(err)
	}

	command, _ := proc.Exe()
	createTime, err := proc.CreateTime()
	if err != nil {
		panic(err)
	}

	return ProcessDetails{
		Command:    command,
		CreateTime: time.Unix(createTime/1000, 0),
	}
}

func getComponentDetails(componentName string) (ProcessDetails, error) {

	ctx := newContext()
	ensureDaemon()
	cl := newClient()
	workspace := requireWorkspace(ctx, cl)
	output, err := workspace.DescribeComponents(ctx, &api.DescribeComponentsInput{})
	if err != nil {
		return ProcessDetails{}, err
	}

	componentFound := false
	state := process.State{}
	for _, component := range output.Components {
		if component.Name == componentName {
			if err := jsonutil.UnmarshalString(component.State, &state); err != nil {
				return ProcessDetails{}, fmt.Errorf("unmarshalling state: %w", err)
			}
			componentFound = true
			break
		}
	}

	if !componentFound {
		return ProcessDetails{}, fmt.Errorf("Component not found")
	}

	if state.Pid == 0 {
		return ProcessDetails{}, fmt.Errorf("Component not running")
	}

	return getProcessDetails(int32(state.Pid)), nil
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Report exo status",
	Long:  `Report exo daemon status.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) > 0 {
			details, err := getComponentDetails(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", details)
			jsonString, err := json.Marshal(details)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonString))
			return nil
		}

		err := loadRunState()
		if err != nil {
			return fmt.Errorf("checking run state: %w", err)
		}

		pid := ""
		if osutil.IsValidPid(runState.Pid) {
			pid = strconv.Itoa(runState.Pid)
		}
		healthy := checkHealthy()

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "healthy:\t%t\n", healthy)
		_, _ = fmt.Fprintf(w, "pid:\t%s\n", pid)
		_, _ = fmt.Fprintf(w, "gui:\t%s\n", runState.URL)
		_ = w.Flush()
		return nil
	},
}

func checkHealthy() bool {
	res, err := http.Get(runState.URL + "_exo/health")
	if err != nil {
		return false
	}
	bs, _ := ioutil.ReadAll(res.Body)
	// See note [HEALTH_CHECK].
	return string(bytes.TrimSpace(bs)) == "ok"
}
