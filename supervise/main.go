package supervise

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/util/cmdutil"
	"github.com/influxdata/go-syslog/v3/rfc5424"
)

var child *os.Process
var varDir string

func Main(command string, args []string) {
	if len(args) < 4 {
		fatalf(`usage: %s <address> <component-id> <working-directory> <program> <args...>

supervise executes and supervises the given command. If successful, the child
pid is written to stdout. The stdout and stderr streams of the supervised process
will be directed to the given address as syslog events.

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

	address := args[0]
	componentID := args[1]
	wd := args[2]
	program := args[3]
	arguments := args[4:]

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fatalf("resolving udp address: %w", err)
	}

	cmd := exec.Command(program, arguments...)
	cmd.Dir = wd
	cmd.Env = os.Environ()

	// Connect pipes.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Start child process.
	if err := cmd.Start(); err != nil {
		fatalf("%v", err)
	}
	child = cmd.Process

	// Reporting child pid to stdout.
	if _, err := fmt.Println(child.Pid); err != nil {
		fatalf("reporting pid: %v", err)
	}

	// Handle signals.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGCHLD)
	go func() {
		for sig := range c {
			switch sig {
			// Forward signals to child.
			case os.Interrupt, os.Kill:
				if err := cmd.Process.Signal(sig); err != nil {
					break
				}
			// Exit when child exits.
			case syscall.SIGCHLD:
				os.Exit(1)
			}
		}
	}()

	// Dial syslog.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fatalf("dialing udp: %w", err)
	}
	defer conn.Close()

	// Proxy logs.
	syslogProcID := strconv.Itoa(child.Pid)
	go pipeToSyslog(ctx, conn, componentID, "out", syslogProcID, stdout)
	go pipeToSyslog(ctx, conn, componentID, "err", syslogProcID, stderr)

	// Wait for child process to exit.
	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
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
			fatalf("sending syslog message: %w", err)
		}
	}
}

const syslogFacility = 1 // "user-level messages".
const syslogSeverity = 6 // "information messages".
const syslogPriority = (syslogFacility * 8) + syslogSeverity

func fatalf(format string, v ...interface{}) {
	if child != nil {
		_ = child.Kill()
	}
	cmdutil.Fatalf(format, v...)
}
