package term

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nerdmaster/terminal"
)

type RawMode struct {
	// Zero value signals stdin.
	FD uintptr

	oldState *terminal.State
}

func (raw *RawMode) Enter() error {
	if raw.oldState != nil {
		return errors.New("already entered raw mode")
	}
	var err error
	raw.oldState, err = terminal.MakeRaw(int(raw.FD))
	return err
}

func (raw *RawMode) Exit() error {
	if raw.oldState == nil {
		return errors.New("terminal state unknown")
	}
	err := terminal.Restore(int(raw.FD), raw.oldState)
	if err == syscall.Errno(0) {
		// terminal.Restore does not properly transform errorno=0 to nil errors.
		err = nil
	}
	if err != nil {
		return err
	}
	raw.oldState = nil
	return nil
}

// Based on NetHack: <https://nethackwiki.com/wiki/%5EZ#cite_note-1>.
func (raw *RawMode) Suspend() error {
	if err := raw.Exit(); err != nil {
		return fmt.Errorf("exiting raw mode: %w", err)
	}
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGCONT)
	if err := syscall.Kill(0, syscall.SIGSTOP); err != nil {
		return fmt.Errorf("signally process to stop: %w", err)
	}
	<-c
	signal.Ignore(syscall.SIGCONT)
	if err := raw.Enter(); err != nil {
		return fmt.Errorf("entering raw mode: %w", err)
	}
	return nil
}
