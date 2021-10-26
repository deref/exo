package exohcl

// Validate deeply analyzes the manifest, collecting all diagnostics eagerly.
func Validate(ctx *AnalysisContext, m *Manifest) {
	m.Analyze(ctx)
	NewEnvironment(m).Analyze(ctx)
	NewComponentSet(m).Analyze(ctx)
}
