package procfile

import (
	"bytes"

	"gopkg.in/alessio/shellescape.v1"
	"mvdan.cc/sh/v3/syntax"
)

func Generate(procs []Process) (string, error) {
	var out bytes.Buffer

	for _, proc := range procs {
		out.WriteString(proc.Name)
		out.WriteString(": ")

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

		syntax.NewPrinter().Print(&out, &syntax.Stmt{
			Cmd: cmd,
		})
		out.WriteByte('\n')
	}

	return out.String(), nil
}

func literalWord(val string) *syntax.Word {
	return &syntax.Word{
		Parts: []syntax.WordPart{
			&syntax.Lit{Value: shellescape.Quote(val)},
		},
	}
}
