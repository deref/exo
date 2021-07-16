module github.com/deref/exo

go 1.16

require (
	github.com/deref/pier v0.0.0-20210620044641-0f71544154e7
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/natefinch/atomic v1.0.1
	github.com/oklog/ulid/v2 v2.0.2
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/spf13/cobra v0.0.5
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

// https://github.com/deref/httpie-go/tree/pier
replace github.com/nojima/httpie-go => github.com/deref/httpie-go v0.7.1-0.20210620034715-00ad2c785a86
