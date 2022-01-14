package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
	machinebox "github.com/machinebox/graphql"
	shurcool "github.com/shurcooL/graphql"
)

type Client struct {
	url        string
	shurcool   *shurcool.Client
	machinebox *machinebox.Client
}

func NewClient(url string, httpClient *http.Client) *Client {
	return &Client{
		url:        url,
		shurcool:   shurcool.NewClient(url, httpClient),
		machinebox: machinebox.NewClient(url, machinebox.WithHTTPClient(httpClient)),
	}
}

// Construct and run a query reflectively from q.
func (cl *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return cl.shurcool.Query(ctx, q, variables)
}

// Construct and run a mutation reflectively from m.
func (cl *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	return cl.shurcool.Mutate(ctx, m, variables)
}

// Run the given query string and decode the response in to resp.
func (cl *Client) Run(ctx context.Context, q string, resp interface{}, vars map[string]interface{}) error {
	req := machinebox.NewRequest(q)
	for k, v := range vars {
		req.Var(k, v)
	}
	return cl.machinebox.Run(ctx, req, resp)
}

func (cl *Client) MutateVoid(ctx context.Context, mutation string, vars map[string]interface{}) error {
	q := FormatVoidMutation(mutation, vars)
	fmt.Println("QUERY:", q)
	var resp struct{}
	return cl.Run(ctx, q, &resp, vars)
}

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

func (cl *Client) StartTask(ctx context.Context, mutation string, vars map[string]interface{}) (jobID string, err error) {
	err = cl.MutateVoid(ctx, mutation, vars)
	return "TODO:JOB_ID", err
}
