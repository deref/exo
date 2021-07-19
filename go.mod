module github.com/deref/exo

go 1.16

require (
	github.com/deref/pier v0.0.0-20210620044641-0f71544154e7
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/fatih/color v1.12.0 // indirect
	github.com/goccy/go-yaml v1.8.10
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/natefinch/atomic v1.0.1
	github.com/oklog/ulid/v2 v2.0.2
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/spf13/cobra v0.0.5
	github.com/zclconf/go-cty v1.8.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

// https://github.com/deref/httpie-go/tree/pier
replace github.com/nojima/httpie-go => github.com/deref/httpie-go v0.7.1-0.20210620034715-00ad2c785a86
