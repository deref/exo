package statefile

type Workspace struct {
	Root       string                `json:"root"`
	Names      map[string]string     `json:"names"`      // Name -> ID.
	Components map[string]*Component `json:"components"` // Keyed by ID.
}

func (ws *Workspace) resolve(refs []string) []*string {
	results := make([]*string, len(refs))
	for i, ref := range refs {
		if _, isID := ws.Components[ref]; isID {
			id := ref
			results[i] = &id
			continue
		}
		id := ws.Names[ref]
		if id != "" {
			results[i] = &id
		}
	}

	return results
}
