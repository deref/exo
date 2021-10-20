package manifest

import "github.com/hashicorp/hcl/v2"

type Diagnostics = hcl.Diagnostics
type Diagnostic = hcl.Diagnostic

const DiagInvalid = hcl.DiagInvalid
const DiagWarning = hcl.DiagWarning
const DiagError = hcl.DiagError

type Pos = hcl.Pos
type Range = hcl.Range
