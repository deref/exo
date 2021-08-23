module github.com/deref/exo

go 1.16

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20210608160410-67692ebc98de
	github.com/BurntSushi/toml v0.3.1
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59
	github.com/containerd/containerd v1.5.4 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/fatih/color v1.12.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/goccy/go-yaml v1.8.10
	github.com/gofrs/flock v0.8.1
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/influxdata/go-syslog/v3 v3.0.0
	github.com/joho/godotenv v1.3.0
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mattn/go-isatty v0.0.13
	github.com/mattn/go-shellwords v1.0.12
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/natefinch/atomic v1.0.1
	github.com/oklog/ulid/v2 v2.0.2
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/shirou/gopsutil/v3 v3.21.6
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/timewasted/go-accept-headers v0.0.0-20130320203746-c78f304b1b09
	github.com/zclconf/go-cty v1.8.0
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
	gopkg.in/alessio/shellescape.v1 v1.0.0-20170105083845-52074bc9df61
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	mvdan.cc/sh/v3 v3.3.1
)

// https://github.com/deref/httpie-go/tree/pier
replace github.com/nojima/httpie-go => github.com/deref/httpie-go v0.7.1-0.20210620034715-00ad2c785a86
