package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/which"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(stateGetCmd)
	stateCmd.AddCommand(stateSetCmd)
	stateCmd.AddCommand(stateClearCmd)
	stateCmd.AddCommand(stateEditCmd)

	stateCmd.AddCommand(helpSubcmd)
}

var stateCmd = &cobra.Command{
	Use:    "state",
	Short:  "View and update the state store.",
	Long:   `Contains subcommands for getting, setting, and clearing state on a per-component basis.`,
	Hidden: true,
	Args:   cobra.NoArgs,
}

var stateGetCmd = &cobra.Command{
	Use:    "get <component>",
	Short:  "Print component state",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()

		workspace := requireWorkspace(ctx, cl)
		output, err := workspace.GetComponentState(ctx, &api.GetComponentStateInput{Ref: componentRef})
		if err != nil {
			return err
		}

		return jsonutil.PrettyPrintJSONString(os.Stdout, output.State)
	},
}

var stateSetCmd = &cobra.Command{
	Use:    "set <component>",
	Short:  "Set component state",
	Long:   "Set component state to the JSON received on stdin.",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()

		var newState map[string]interface{}
		if err := json.NewDecoder(os.Stdin).Decode(&newState); err != nil {
			return fmt.Errorf("reading state from stdin: %w", err)
		}

		workspace := requireWorkspace(ctx, cl)
		_, err := workspace.SetComponentState(ctx, &api.SetComponentStateInput{
			Ref:   componentRef,
			State: jsonutil.MustMarshalString(newState),
		})
		return err
	},
}

var stateClearCmd = &cobra.Command{
	Use:    "clear <component>",
	Short:  "Clear component state",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()

		workspace := requireWorkspace(ctx, cl)
		_, err := workspace.SetComponentState(ctx, &api.SetComponentStateInput{
			Ref:   componentRef,
			State: "{}",
		})
		return err
	},
}

var stateEditCmd = &cobra.Command{
	Use:    "edit <component>",
	Short:  "Edit component state",
	Long:   "Edit component state using your preferred editor.",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()

		workspace := requireWorkspace(ctx, cl)
		output, err := workspace.GetComponentState(ctx, &api.GetComponentStateInput{Ref: componentRef})
		if err != nil {
			return err
		}

		tmpfile, err := ioutil.TempFile("", "state.*.json")
		if err != nil {
			return fmt.Errorf("creating temporary file: %w", err)
		}
		defer os.Remove(tmpfile.Name())

		if err := jsonutil.PrettyPrintJSONString(tmpfile, output.State); err != nil {
			return fmt.Errorf("writing to temporary file: %w", err)
		}
		if err := tmpfile.Close(); err != nil {
			return fmt.Errorf("closing temporary file: %w", err)
		}

		stat, err := os.Stat(tmpfile.Name())
		if err != nil {
			return fmt.Errorf("checking modification time: %w", err)
		}
		originalModTime := stat.ModTime()

		editor := os.Getenv("EDITOR")
		if editor == "" {

			for _, candidateEditor := range []string{
				"sensible-editor",
				"editor",
				"code",
				"vim",
				"nano",
				"vi",
				"emacs",
				"ee",
			} {
				found, err := which.Which(candidateEditor)
				if err != nil && !strings.Contains(err.Error(), "not found") {
					return fmt.Errorf("looking up candidate editor %q: %w", candidateEditor, err)
				}
				if found != "" {
					editor = found
					break
				}
			}

			if editor == "" {
				return errors.New("No editor available")
			}
		}

		edit := exec.Command(editor, tmpfile.Name())
		edit.Stdin = os.Stdin
		edit.Stdout = os.Stdout
		edit.Stderr = os.Stderr
		if err := edit.Run(); err != nil {
			return fmt.Errorf("editing state file: %w", err)
		}

		if stat, err = os.Stat(tmpfile.Name()); err != nil {
			return fmt.Errorf("checking modification time: %w", err)
		}
		newModTime := stat.ModTime()
		if !newModTime.After(originalModTime) {
			fmt.Fprintf(os.Stderr, "Not modified - not updating state.")
			return nil
		}

		newState, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return fmt.Errorf("reading updated state: %w", err)
		}

		_, err = workspace.SetComponentState(ctx, &api.SetComponentStateInput{
			Ref:   componentRef,
			State: string(newState),
		})
		return err
	},
}
