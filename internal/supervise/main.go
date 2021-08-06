package supervise

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/logd/api"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/sysutil"
	"github.com/influxdata/go-syslog/v3/rfc5424"
)

var child *os.Process
var varDir string

func Main(command string, args []string) {
	if len(args) < 4 {
		fatalf(`usage: %s <syslog-port> <component-id> <working-directory> <timeout> <env> <program> <args...>

supervise executes and supervises the given command. If successful, the child
pid is written to stdout. The stdout and stderr streams of the supervised process
will be directed to the given port on localhost as syslog events.

Syslog messages use the following fields:

APPNAME = Component ID for the Exo process that is being logged.
PROCID = PID of the supervised process. As per RFC5425, this field has "no
				 interoperable meaning, except that a change in the value indicates
				 there has been a discontinuity in syslog reporting".
MSGID = The message "type". Set to "out" or "err" to specify which stdio
				stream the message came from.
`, command)
	}
	ctx := context.Background()

	syslogPort := args[0]
	componentID := args[1]
	wd := args[2]
	timeout := args[3]
	envString := args[4]
	program := args[5]
	arguments := args[6:]

	var crashFile *os.File
	cleanExit := func(statusCode int) {
		if crashFile != nil {
			_ = os.Remove(crashFile.Name())
		}
		os.Exit(statusCode)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:"+syslogPort)
	if err != nil {
		fatalf("resolving udp address: %w", err)
	}

	timeoutSeconds, timeoutErr := strconv.Atoi(timeout)
	if timeoutErr != nil {
		fatalf(timeoutErr.Error())
	}

	childEnv := make(map[string]string)
	if err := json.Unmarshal([]byte(envString), &childEnv); err != nil {
		fatalf("decoding environment variables from %q: %v", envString, err)
	}

	cmd := exec.Command(program, arguments...)
	cmd.Dir = wd
	cmd.Env = os.Environ()
	for key, val := range childEnv {
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

	hasSignalledChildToQuit := false
	// Handle signals.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGCHLD)
	go func() {
		for sig := range c {
			switch sig {
			// Forward signals to child.
			case os.Interrupt, syscall.SIGTERM:
				hasSignalledChildToQuit = true
				if err := cmd.Process.Signal(sig); err != nil {
					break
				}
				// After some timeout kill the entire process group.
				time.Sleep(time.Second * time.Duration(timeoutSeconds))
				die()

			// Exit when child exits.
			case syscall.SIGCHLD:
				if hasSignalledChildToQuit {
					cleanExit(0)
				}
				cleanExit(1)
			}
		}
	}()

	// Start child process.
	if err := cmd.Start(); err != nil {
		fatalf("%v", err)
	}
	child = cmd.Process

	// Dial syslog.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fatalf("dialing udp: %v", err)
	}
	defer conn.Close()

	// Reporting child pid to stdout.
	if _, err := fmt.Println(child.Pid); err != nil {
		fatalf("reporting pid: %v", err)
	}

	// NOTE [SUPERVISE_STDERR]: The "started ok" message will release and readers
	// who are waiting for a message on stderr. Then we redirect stderr to a temp
	// file so that if any supervision failures happen, we have a crash log we
	// can inspect.
	_, _ = fmt.Fprintf(os.Stderr, "started ok\n")
	crashFile, _ = ioutil.TempFile("", "supervise.*.stderr")
	if crashFile != nil {
		_ = sysutil.Dup2(int(crashFile.Fd()), 2)
		fmt.Fprintf(os.Stderr, "supervisor pid: %d\n", os.Getpid())
	}

	// Proxy logs.
	syslogProcID := strconv.Itoa(child.Pid)
	go pipeToSyslog(ctx, conn, componentID, "out", syslogProcID, stdout)
	go pipeToSyslog(ctx, conn, componentID, "err", syslogProcID, stderr)

	// Wait for child process to exit.
	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		cleanExit(exitErr.ExitCode())
	}
	if err != nil {
		fatalf("wait error: %v", err)
	}
}

func pipeToSyslog(ctx context.Context, conn net.Conn, componentID string, name string, procID string, r io.Reader) {
	b := bufio.NewReaderSize(r, api.MaxMessageSize)
	for {
		message, isPrefix, err := b.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			fatalf("reading %s: %v", name, err)
		}
		// TODO: Do something better with lines that are too long.
		for isPrefix {
			// Skip remainder of line.
			message = append([]byte{}, message...)
			_, isPrefix, err = b.ReadLine()
			if err == io.EOF {
				return
			}
			if err != nil {
				fatalf("reading %s: %v", name, err)
			}
		}

		sm := &rfc5424.SyslogMessage{}
		sm.SetVersion(1)
		sm.SetPriority(syslogPriority)
		sm.SetTimestamp(time.Now().Format(chrono.RFC3339MicroUTC))
		sm.SetAppname(componentID)
		sm.SetProcID(procID)
		sm.SetMsgID(name) // See note: [LOG_COMPONENTS].
		sm.SetMessage(string(message))
		packet, err := sm.String()
		if err != nil {
			fatalf("building syslog message: %w", err)
		}
		if _, err := io.WriteString(conn, packet); err != nil {
			log.Printf("sending syslog message: %v", err)
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
	pgrp := syscall.Getpgrp()
	_ = osutil.KillProcessGroup(pgrp)
}
