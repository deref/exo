package cli

import (
	"bytes"
	"html/template"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/spf13/cobra"
)

func controlComponents(cmd *cobra.Command, args []string, workspaceMutation string, componentsMutation string, vars map[string]interface{}) error {
	ctx := cmd.Context()
	checkOrEnsureServer()
	kernel := newClient().Kernel()

	cl, shutdown := dialGraphQL(ctx)
	defer shutdown()

	// TODO: It would be nice to have generated mutation methods.
	var tmpl *template.Template
	var data struct {
		Mutation string
	}
	vars = jsonutil.Merge(map[string]interface{}{
		"workspace": currentWorkspaceRef(),
	}, vars)
	if len(args) == 0 {
		tmpl = template.Must(template.New("").Parse(`
			mutation (
				$workspace: String!
			) {
				{{ .Mutation }}(workspace: $workspace) {
					id
				}
			}
		`))
		data.Mutation = workspaceMutation
	} else {
		tmpl = template.Must(template.New("").Parse(`
			mutation (
				$workspace: String!
				$components: [String!]!
			) {
				{{ .Mutation }}(workspace: $workspace, components: $components) {
					id
				}
			}
		`))
		data.Mutation = componentsMutation
		vars["components"] = args
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	var m struct {
		Job struct {
			ID string
		}
	}
	if err := cl.Run(ctx, buf.String(), &m, vars); err != nil {
		return err
	}

	return watchJob(ctx, kernel, m.Job.ID)
}
