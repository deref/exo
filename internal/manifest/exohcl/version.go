package exohcl

import "fmt"

var Latest = FormatVersion{
	Major: 1,
	Minor: 0,
}

var Starter = fmt.Sprintf(`# See https://docs.deref.io/exo for details.

exo = "%s"
`, Latest)

type FormatVersion struct {
	Major int
	Minor int
}

func (ver FormatVersion) String() string {
	return fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
}
