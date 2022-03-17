package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(stateGetCmd)
	stateCmd.AddCommand(stateSetCmd)
	stateCmd.AddCommand(stateClearCmd)
	stateCmd.AddCommand(stateEditCmd)

	stateCmd.AddCommand(makeHelpSubcmd())
}

var stateCmd = &cobra.Command{
	Use:    "state",
	Short:  "View and update the state store.",
	Long:   `Contains subcommands for getting, setting, and clearing state on a per-component basis.`,
	Hidden: true,
	Args:   cobra.NoArgs,
}

var stateGetCmd = &cobra.Command{
	Use:   "get <component>",
	Short: "Print component state",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := cmd.Context()
		cl := newClient()

		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.GetComponentState(ctx, &api.GetComponentStateInput{Ref: componentRef})
		if err != nil {
			return err
		}

		return jsonutil.PrettyPrintJSONString(os.Stdout, output.State)
	},
}

var stateSetCmd = &cobra.Command{
	Use:   "set <component>",
	Short: "Set component state",
	Long:  "Set component state to the JSON received on stdin.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := cmd.Context()
		cl := newClient()

		var newState map[string]any
		if err := json.NewDecoder(os.Stdin).Decode(&newState); err != nil {
			return fmt.Errorf("reading state from stdin: %w", err)
		}

		workspace := requireCurrentWorkspace(ctx, cl)
		_, err := workspace.SetComponentState(ctx, &api.SetComponentStateInput{
			Ref:   componentRef,
			State: jsonutil.MustMarshalString(newState),
		})
		return err
	},
}

var stateClearCmd = &cobra.Command{
	Use:   "clear <component>",
	Short: "Clear component state",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := cmd.Context()
		cl := newClient()

		workspace := requireCurrentWorkspace(ctx, cl)
		_, err := workspace.SetComponentState(ctx, &api.SetComponentStateInput{
			Ref:   componentRef,
			State: "{}",
		})
		return err
	},
}

var stateEditCmd = &cobra.Command{
	Use:   "edit <component>",
	Short: "Edit component state",
	Long:  "Edit component state using your preferred editor.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentRef := args[0]
		ctx := cmd.Context()
		cl := newClient()

		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.GetComponentState(ctx, &api.GetComponentStateInput{Ref: componentRef})
		if err != nil {
			return err
		}

		// TODO: pretty-print / minify json.
		oldState := output.State
		newState, err := term.EditString("state.*.json", oldState)

		_, err = workspace.SetComponentState(ctx, &api.SetComponentStateInput{
			Ref:   componentRef,
			State: string(newState),
		})
		return err
	},
}
