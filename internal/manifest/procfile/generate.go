package procfile

import (
	"fmt"
	"io"

	"gopkg.in/alessio/shellescape.v1"
	"mvdan.cc/sh/v3/syntax"
)

func Generate(w io.Writer, procs []Process) error {
	for _, proc := range procs {
		fmt.Fprintf(w, "%s: ", proc.Name)

		cmd := &syntax.CallExpr{
			Args: []*syntax.Word{
				literalWord(proc.Program),
			},
		}
		for _, arg := range proc.Arguments {
			cmd.Args = append(cmd.Args, literalWord(arg))
		}

		for key, val := range proc.Environment {
			cmd.Assigns = append(cmd.Assigns, &syntax.Assign{
				Name:  &syntax.Lit{Value: key},
				Value: literalWord(val),
			})
		}

		syntax.NewPrinter().Print(w, &syntax.Stmt{
			Cmd: cmd,
		})
		fmt.Fprint(w, "\n")
	}

	return nil
}

func literalWord(val string) *syntax.Word {
	return &syntax.Word{
		Parts: []syntax.WordPart{
			&syntax.Lit{Value: shellescape.Quote(val)},
		},
	}
}
