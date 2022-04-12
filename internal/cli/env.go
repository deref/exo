package cli

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.Flags().StringVar(&envCmdFlags.Scope, "scope", "workspace", "workspace, cluster, stack, or component")
	// TODO: Flags for getting environment of specific scopes by ref or
	// components in stacks besides the current one.
	envCmd.Flags().StringVar(&envCmdFlags.Component, "component", "", "component ref, implies --scope=component")
	// TODO: Handle showing/hiding of sensitive values.
}

var envCmdFlags struct {
	Scope     string
	Component string
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long: `Prints environment variables in .env format.

The workspace's environment is that of the current stack. Or, if there is
no current stack, it's the environment of the local cluster`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		scope := envCmdFlags.Scope
		if envCmdFlags.Component != "" {
			if !cmd.Flags().Changed("scope") {
				scope = "component"
			} else if scope != "component" {
				return fmt.Errorf("--component conflicts with --scope=%q", scope)
			}
		}
		var environment environmentFragment
		var err error
		switch scope {
		case "cluster":
			var q struct {
				Cluster *struct {
					Environment environmentFragment
				} `graphql:"defaultCluster"`
			}
			err = api.Query(ctx, svc, &q, nil)
		case "stack":
			var q struct {
				Stack *struct {
					Environment environmentFragment
				} `graphql:"stackByRef(ref: $stack)"`
			}
			err = api.Query(ctx, svc, &q, map[string]any{
				"stack": currentStackRef(),
			})
		case "workspace":
			environment = getWorkspaceEnvironment(ctx)
		case "component":
			var q struct {
				Stack *struct {
					Environment environmentFragment
				} `graphql:"stackByRef(ref: $stack)"`
			}
			err = api.Query(ctx, svc, &q, map[string]any{
				"stack":     currentStackRef(),
				"component": envCmdFlags.Component,
			})
		default:
			return fmt.Errorf("unknown scope: %q", scope)
		}
		if err != nil {
			return err
		}
		showEnvironment(environment)
		return nil
	},
}

func getWorkspaceEnvironment(ctx context.Context) environmentFragment {
	var q struct {
		Workspace *struct {
			Environment environmentFragment
		} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
	}
	mustQueryWorkspace(ctx, &q, nil)
	return q.Workspace.Environment
}

func showEnvironment(env environmentFragment) {
	for _, v := range env.Variables {
		fmt.Println(osutil.FormatDotEnvEntry(v.Name, v.Value))
	}
}

type environmentFragment struct {
	Variables []struct {
		Name  string
		Value string
	}
}

func environmentFragmentToMap(fragment environmentFragment) map[string]string {
	variables := fragment.Variables
	m := make(map[string]string, len(variables))
	for _, variable := range variables {
		m[variable.Name] = variable.Value
	}
	return m
}
