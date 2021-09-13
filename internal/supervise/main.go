package supervise

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/logd/api"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/sysutil"
	"github.com/influxdata/go-syslog/v3/rfc5424"
)

var varDir string
var pgrp int

func Main() {
	var crashFile *os.File
	cleanExit := func() {
		if crashFile != nil {
			_ = os.Remove(crashFile.Name())
		}
		os.Exit(0)
	}

	pgrp = syscall.Getpgrp()
	ctx := context.Background()

	cfg := &Config{}
	if err := json.NewDecoder(os.Stdin).Decode(cfg); err != nil {
		fatalf("reading config from stdin: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		fatalf("validating config: %v", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", cfg.SyslogPort))
	if err != nil {
		fatalf("resolving udp address: %v", err)
	}

	cmd := exec.Command(cfg.Program, cfg.Arguments...)
	cmd.Dir = cfg.WorkingDirectory
	cmd.Env = os.Environ() // TODO: Should we start with an empty env?
	for key, val := range cfg.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
	}

	// Connect pipes.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Dial syslog.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fatalf("dialing udp: %v", err)
	}
	defer conn.Close()

	// Register for signals.  Do this before starting the child to
	// guarantee we see any exist of a child process.
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGCHLD)

	// Start child process.
	if err := cmd.Start(); err != nil {
		fatalf("%v", err)
	}
	child := cmd.Process

	// Reporting child pid to stdout.
	if _, err := fmt.Println(child.Pid); err != nil {
		fatalf("reporting pid: %v", err)
	}

	// NOTE [SUPERVISE_STDERR]: The "started ok" message will release any readers
	// who are waiting for a message on stderr. Then we redirect stderr to a temp
	// file so that if any supervision failures happen, we have a crash log we
	// can inspect.
	_, _ = fmt.Fprintf(os.Stderr, "started ok\n")
	crashFile, _ = ioutil.TempFile("", "supervise.*.stderr")
	if crashFile != nil {
		_ = sysutil.Dup2(int(crashFile.Fd()), 2)
	}

	log.Println("supervisor pid:", os.Getpid())
	log.Println("child pid:", child.Pid)

	// Asynchronously handle SIGCHILD. We spawn this goroutine after starting the
	// child, so that we don't have to coordinate concurrent access to the child
	// variable.
	go func() {
		for sig := range c {
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				// We expect exo to send these to the whole group. This means that a
				// well behaved child will handle SIGTERM, leading to us receiving
				// SIGCHLD. However, we must ignore these signals so that we don't stop
				// processing logs before the child stops sending them!
			case syscall.SIGCHLD:
				// Allow a little extra time to gather shutdown logs from the child.
				<-time.After(1 * time.Second)
				cleanExit()
			}
		}
	}()

	// Proxy logs.
	syslogProcID := strconv.Itoa(child.Pid)

	var wg sync.WaitGroup
	work := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	work(func() {
		pipeToSyslog(ctx, conn, cfg.ComponentID, "out", syslogProcID, stdout)
	})
	work(func() {
		pipeToSyslog(ctx, conn, cfg.ComponentID, "err", syslogProcID, stderr)
	})

	// Wait for child process and log forwarding to exit.
	err = cmd.Wait()
	wg.Wait()
	if _, ok := err.(*exec.ExitError); ok {
		cleanExit()
	}
	if err != nil {
		fatalf("wait error: %v", err)
	}
}

func pipeToSyslog(ctx context.Context, conn net.Conn, componentID string, name string, procID string, r io.Reader) {
	b := bufio.NewReaderSize(r, api.MaxMessageSize)
	readLine := func() (string, error) {
		// Usage of ReadLine in preference to ReadString is intentional, since
		// ReadString will perform unbounded buffering.
		// See discussion here: https://github.com/deref/exo/pull/322
		message, isPrefix, err := b.ReadLine()
		for isPrefix {
			// Skip remainder of line.
			message = append([]byte{}, message...)
			_, isPrefix, err = b.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				fatalf("reading %s: %v", name, err)
			}
		}
		return string(message), err
	}

	for {
		message, err := readLine()

		// Error handling is performed after piping the message to syslog since we
		// always want to write the message, even if an error has occurred.
		if message != "" {
			if message[len(message)-1] == '\n' {
				message = message[:len(message)-1]
			}
			sm := &rfc5424.SyslogMessage{}
			sm.SetVersion(1)
			sm.SetPriority(syslogPriority)
			sm.SetTimestamp(chrono.Now(ctx).Format(chrono.RFC3339MicroUTC))
			sm.SetAppname(componentID)
			sm.SetProcID(procID)
			sm.SetMsgID(name) // See note: [LOG_COMPONENTS].
			sm.SetMessage(message)
			packet, err := sm.String()
			if err != nil {
				fatalf("building syslog message: %w", err)
			}
			if _, err := io.WriteString(conn, packet); err != nil {
				log.Printf("sending syslog message: %v", err)
			}
		}
		if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
			return
		}
		if err != nil {
			fatalf("reading %s: %v", name, err)
		}
	}
}

const syslogFacility = 1 // "user-level messages".
const syslogSeverity = 6 // "information messages".
const syslogPriority = (syslogFacility * 8) + syslogSeverity

func fatalf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	die()
}

func die() {
	_ = osutil.KillGroup(pgrp)
}
