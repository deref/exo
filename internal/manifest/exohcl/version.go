package exohcl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

const LatestMajor = 0
const LatestMinor = 1

var Latest = fmt.Sprintf("%d.%d", LatestMajor, LatestMinor)

var Starter = fmt.Sprintf(`# See https://docs.deref.io/exo for details.

exo = "%s"
`, Latest)

type FormatVersion struct {
	// Analysis inputs.
	Attribute *hcl.Attribute

	// Analysis outputs.
	Major int
	Minor int
}

func NewFormatVersion(m *Manifest) *FormatVersion {
	return &FormatVersion{
		Attribute: m.Content.Attributes["exo"],
	}
}

func (ver *FormatVersion) Analyze(ctx *AnalysisContext) {
	if ver.Attribute == nil {
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "exo version attribute is required",
		})
		return
	}
	s, diag := parseLiteralString(ver.Attribute.Expr)
	if diag != nil {
		ctx.AppendDiags(diag)
		return
	}
	parts := strings.Split(s, ".")
	ok := len(parts) == 2
	ints := make([]int, len(parts))
	for index, part := range parts {
		i, err := strconv.Atoi(part)
		ok = ok && err == nil && i >= 0
		ints[index] = i
	}
	if !ok {
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid exo format version constraint",
			Detail:   fmt.Sprintf(`Exo format constraint should be specified as 'major.minor'. The latest version is %q.`, Latest),
			Subject:  ver.Attribute.Expr.Range().Ptr(),
			Context:  ver.Attribute.Range.Ptr(),
		})
		return
	}
	ver.Major = ints[0]
	ver.Minor = ints[1]
	switch ver.Major {
	case 0:
		if LatestMinor < ver.Minor {
			ctx.AppendDiags(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported exo format minor version",
				Detail:   fmt.Sprintf(`Unsupported exo format minor version. The maximum supported "0.x" version is %q.`, Latest),
				Subject:  ver.Attribute.Expr.Range().Ptr(),
				Context:  ver.Attribute.Range.Ptr(),
			})
		}
	default:
		ctx.AppendDiags(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported exo format major version",
			Detail:   fmt.Sprintf(`Unsupported exo format major version. The latest version is %q.`, Latest),
			Subject:  ver.Attribute.Expr.Range().Ptr(),
			Context:  ver.Attribute.Range.Ptr(),
		})
	}
	return
}
