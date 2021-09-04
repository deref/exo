package procfile

import (
	"fmt"
	"sort"
	"strconv"
)

// Establishes a stable order of Processes based on PORT environment variable
// assignment logic. Removes redundant PORT variables from environments.  If
// gaps are left by the port alignment mapping, they are filled in with the
// remaining processes sorted alphabetically by name.
func Organize(in *[]Process) {
	n := len(*in)
	out := make([]Process, n)
	nameToProcess := make(map[string]Process)
	extras := make([]string, 0, n)
	for i, proc := range *in {
		// Copy port-aligned processes to output, or putin to extras list.
		var port int
		if proc.Environment != nil {
			port, _ = strconv.Atoi(proc.Environment["PORT"])
		}
		index := (port - BasePort) / PortStep
		remainder := (port - BasePort) % PortStep
		if 0 <= index && index < n && remainder == 0 {
			delete(proc.Environment, "PORT")
			out[index] = proc
		} else {
			extras = append(extras, proc.Name)
		}

		// Store process by name.
		if proc.Name == "" {
			panic(fmt.Errorf("unnamed process at index %d", i))
		}
		nameToProcess[proc.Name] = proc
	}

	// Sort extras list.
	sort.Strings(extras)

	// Copy non-port aligned processes to fill in gaps.
	for i, proc := range out {
		if proc.Name == "" {
			out[i] = nameToProcess[extras[0]]
			extras = extras[1:]
		}
	}

	*in = out
}
