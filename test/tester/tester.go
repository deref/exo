package tester

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ExoTest struct {
	Test       func(ctx context.Context, t ExoTester) error
	FixtureDir string
}

type ExoTester struct {
	serverPort int
	guiPort    int
	exoHome    string
	exoBinary  string
	fixtureDir string
	logger     *logrus.Logger
	logBuffer  io.Reader
}

func (et ExoTester) GetExoLogs() (string, error) {
	mainLogs, err := ioutil.ReadFile(filepath.Join(et.exoHome, "var", "exod.log"))
	if err != nil {
		return "", fmt.Errorf("reading exod.log: %w", err)
	}
	stdoutLogs, err := ioutil.ReadFile(filepath.Join(et.exoHome, "var", "exod.stdout"))
	if err != nil {
		return "", fmt.Errorf("reading exod.stdout: %w", err)
	}
	stderrLogs, err := ioutil.ReadFile(filepath.Join(et.exoHome, "var", "exod.stderr"))
	if err != nil {
		return "", fmt.Errorf("reading exod.stderr: %w", err)
	}
	workspaceLogs, _, err := et.RunExo(context.TODO(), "logs --no-follow")
	if err != nil {
		et.logger.Warn("Failed to get workspace logs: ", err)
	}
	return fmt.Sprintf("Daemon logs:\n%s\nStdout logs:\n%s\nStderr logs:\n%s\nWorkspace logs:\n%s\n", mainLogs, stdoutLogs, stderrLogs, workspaceLogs), nil
}

func (et ExoTester) RunTest(ctx context.Context, test ExoTest) (io.Reader, error) {
	defer et.StopDaemon(context.Background())
	if err := et.StartDaemon(ctx); err != nil {
		return et.logBuffer, fmt.Errorf("starting daemon: %w", err)
	}
	et.logger.Debug("Started exo")
	if err := test.Test(ctx, et); err != nil {
		return et.logBuffer, fmt.Errorf("running test: %w", err)
	}
	et.logger.Debug("Finished test")
	if err := et.StopDaemon(ctx); err != nil {
		return et.logBuffer, fmt.Errorf("stopping daemon: %w", err)
	}
	return et.logBuffer, nil
}

func (et ExoTester) WaitTillProcessesReachState(ctx context.Context, state string, names []string) error {
	errTimeout := fmt.Errorf("timed out waiting for %q to reach %s", strings.Join(names, ", "), state)
	for {
		processes, err := et.PS(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return errTimeout
			default:
				return err
			}
		}

		doneProcesses := []string{}
		for _, proc := range processes {
			if proc.Status == state {
				doneProcesses = append(doneProcesses, proc.Name)
			}
		}

		sort.Strings(names)
		sort.Strings(doneProcesses)
		if len(names) == len(doneProcesses) && reflect.DeepEqual(names, doneProcesses) {
			return nil
		}

		select {
		case <-ctx.Done():
			return errTimeout
		case <-time.After(100 * time.Millisecond):
		}
	}
}

// Runs an exo CLI command and blocks until it terminates.
func (et ExoTester) RunExo(ctx context.Context, arguments ...string) (stdout, stderr string, err error) {
	path, _ := os.LookupEnv("PATH")
	cmd := exec.CommandContext(ctx, et.exoBinary, arguments...)
	cmd.Dir = et.fixtureDir
	cmd.Env = append(cmd.Env, "EXO_HOME="+et.exoHome)
	cmd.Env = append(cmd.Env, "PATH="+path)
	var stdoutBuffer, stderrBuffer bytes.Buffer
	logWriter := et.logger.Writer()
	defer logWriter.Close()
	stdoutWriter := io.MultiWriter(&stdoutBuffer, logWriter)
	stderrWriter := io.MultiWriter(&stderrBuffer, logWriter)
	cmd.Stderr = stderrWriter
	cmd.Stdout = stdoutWriter
	et.logger.Debug("Running exo ", cmd.Args)
	et.logger.Debug("env ", cmd.Env)
	et.logger.Debug("wd ", cmd.Dir)
	err = cmd.Run()
	et.logger.Debug("exit code: ", cmd.ProcessState.ExitCode())
	return stdoutBuffer.String(), stderrBuffer.String(), err
}

type psProcessInfo struct {
	Name   string
	ID     string
	Status string
	Kind   string
}

func (et ExoTester) PS(ctx context.Context) ([]psProcessInfo, error) {
	stdout, _, err := et.RunExo(ctx, "ps")
	if err != nil {
		return nil, fmt.Errorf("running exo ps: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	processes := make([]psProcessInfo, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("invalid exo ps output: %q", line)
		}
		processes[i] = psProcessInfo{
			Name:   fields[0],
			ID:     fields[1],
			Status: fields[2],
			Kind:   fields[3],
		}
	}

	return processes, nil
}

type statusResult struct {
	Healthy bool
	PID     int
	GUI     string
}

func (et ExoTester) GetStatus(ctx context.Context) (statusResult, error) {
	status, _, err := et.RunExo(ctx, "status")
	if err != nil {
		return statusResult{}, err
	}

	lines := strings.Split(strings.TrimSpace(status), "\n")
	if len(lines) != 3 {
		return statusResult{}, fmt.Errorf("invalid status result: %q", status)
	}

	pid := 0
	pidFields := strings.Fields(lines[1])
	if len(pidFields) > 1 {
		pid, err = strconv.Atoi(pidFields[1])
		if err != nil {
			return statusResult{}, fmt.Errorf("converting pid (%q) to int: %w", pidFields[1], err)
		}
	}

	return statusResult{
		Healthy: strings.Fields(lines[0])[1] == "true",
		PID:     pid,
		GUI:     strings.Fields(lines[2])[1],
	}, nil
}

func (et ExoTester) StartDaemon(ctx context.Context) error {
	if _, _, err := et.RunExo(ctx, "daemon"); err != nil {
		return err
	}
	timeout := time.Now().Add(time.Second * 30)
	for time.Now().Before(timeout) {
		status, _ := et.GetStatus(ctx)
		if status.Healthy {
			return nil
		}
	}
	return errors.New("timed out waiting for daemon to be healthy")
}

func (et ExoTester) StopDaemon(ctx context.Context) error {
	// FIXME: This explicit "exo stop" is necessary because right now stopping the
	// daemon doesn't appear to stop the underlying docker images. Killing the
	// daemon should probably kill all attached docker images.
	_, _, err := et.RunExo(ctx, "stop", "--timeout=0")
	if err != nil {
		return fmt.Errorf("running exo stop: %w", err)
	}

	_, _, err = et.RunExo(ctx, "exit")
	if err != nil {
		return fmt.Errorf("running exo exit: %w", err)
	}
	status, err := et.GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("getting exo status: %w", err)
	}
	if status.Healthy {
		return fmt.Errorf("failed to shutdown exo")
	}
	return err
}

func mustGetFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

type testConfig struct {
	configContents string
	serverPort     int
	guiPort        int
	syslogPort     int
}

func getTestConfig() testConfig {
	serverPort := mustGetFreePort()
	guiPort := mustGetFreePort()
	syslogPort := mustGetFreePort()
	configContents := fmt.Sprintf(`
httpPort = %d
[client]
url = "http://localhost:%d"

[gui]
port = %d

[log]
syslogPort = %d

[telemetry]
derefInternalUser = true
`, serverPort, serverPort, guiPort, syslogPort)

	return testConfig{configContents: configContents,
		guiPort:    guiPort,
		serverPort: serverPort,
		syslogPort: syslogPort,
	}
}

func MakeExoTester(exoBinPath, fixtureBasePath string, test ExoTest) ExoTester {
	config := getTestConfig()

	logBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
	logger.SetOutput(logBuffer)
	logger.SetLevel(logrus.DebugLevel)

	exoBinPath, _ = filepath.Abs(exoBinPath)
	if _, err := os.Stat(exoBinPath); err != nil {
		fmt.Println("Cannot stat exo binary:", err)
		os.Exit(1)
	}

	fixtureDir := filepath.Join(fixtureBasePath, test.FixtureDir)
	if _, err := os.Stat(fixtureDir); err != nil {
		fmt.Println("Could not read fixture directory:", err)
		os.Exit(1)
	}

	exoHome, err := os.MkdirTemp("", "exo-temp-home-")
	if err != nil {
		fmt.Println("Could not create temp home:", err)
		os.Exit(1)
	}

	logger.Infof("EXO_HOME is %q", exoHome)

	err = ioutil.WriteFile(filepath.Join(exoHome, "config.toml"), []byte(config.configContents), 0600)
	if err != nil {
		fmt.Println("could not write config file:", err)
		os.Exit(1)
	}

	tester := ExoTester{
		serverPort: config.serverPort,
		guiPort:    config.guiPort,
		exoBinary:  exoBinPath,
		fixtureDir: fixtureDir,
		exoHome:    exoHome,
		logger:     logger,
		logBuffer:  logBuffer,
	}
	return tester
}
