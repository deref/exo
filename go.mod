module github.com/deref/exo

go 1.16

require (
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59
	github.com/containerd/containerd v1.5.4 // indirect
	github.com/deref/pier v0.0.0-20210620044641-0f71544154e7
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/goccy/go-yaml v1.8.10
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/influxdata/go-syslog/v3 v3.0.0
	github.com/mattn/go-isatty v0.0.13
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/natefinch/atomic v1.0.1
	github.com/oklog/ulid/v2 v2.0.2
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/timewasted/go-accept-headers v0.0.0-20130320203746-c78f304b1b09
	github.com/zclconf/go-cty v1.8.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

// https://github.com/deref/httpie-go/tree/pier
replace github.com/nojima/httpie-go => github.com/deref/httpie-go v0.7.1-0.20210620034715-00ad2c785a86
