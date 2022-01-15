package api

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
)

// XXX use stuff in this file to work tasks.

func FormatVoidMutation(mutation string, vars map[string]interface{}) string {
	arguments := make([]*ast.Argument, 0, len(vars))
	for k, v := range vars {
		arguments = append(arguments, &ast.Argument{
			Kind:  kinds.Argument,
			Name:  newNameNode(k),
			Value: newValueNode(v),
		})
	}
	doc := &ast.Document{
		Kind: kinds.Document,
		Definitions: []ast.Node{
			&ast.OperationDefinition{
				Kind:      kinds.OperationDefinition,
				Operation: "mutation",
				SelectionSet: &ast.SelectionSet{
					Kind: kinds.SelectionSet,
					Selections: []ast.Selection{
						&ast.Field{
							Kind: kinds.Field,
							Name: &ast.Name{
								Kind:  kinds.Name,
								Value: mutation,
							},
							Arguments: arguments,
							SelectionSet: &ast.SelectionSet{
								Kind: kinds.SelectionSet,
								Selections: []ast.Selection{
									&ast.Field{
										Kind: kinds.Field,
										Name: newNameNode("__typename"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return printer.Print(doc).(string)
}

func newNameNode(value string) *ast.Name {
	return &ast.Name{
		Kind:  kinds.Name,
		Value: value,
	}
}

func newValueNode(value interface{}) ast.Value {
	switch value := value.(type) {
	case string:
		return &ast.StringValue{
			Kind:  kinds.StringValue,
			Value: value,
		}
	default:
		panic(fmt.Errorf("cannot convert %T to GraphQL ast node", value))
	}
}
