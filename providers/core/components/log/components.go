package log

// NOTE [LOG_COMPONENTS]: We don't yet treat logs as components of their own,
// so we hard code an expansion from process -> stdout/stderr log pairs.
// Multiple places in the code make brittle assumptions about this and are
// tagged with this note accordingly.
func ComponentLogNames(provider string, componentID string) []string {
	switch provider {
	case "unix":
		return []string{
			componentID + ":out",
			componentID + ":err",
		}
	case "docker":
		return []string{componentID}
	default:
		return nil
	}
}
