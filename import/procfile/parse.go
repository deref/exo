package procfile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Procfile struct {
	Processes []Process
}

type Process struct {
	Name      string
	Command   string
	Arguments []string
}

func Parse(r io.Reader) (*Procfile, error) {
	var procfile Procfile
	br := bufio.NewReader(r)
	lineIndex := 0
	for {
		lineIndex++
		line, isPrefix, err := br.ReadLine()
		if io.EOF == err {
			break
		}
		if isPrefix {
			return nil, fmt.Errorf("line %d is too long", lineIndex)
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
			// Blank or comment.
			continue
		}
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d is invalid", lineIndex)
		}
		name := string(parts[0])
		// TODO: Validate name is alphanumeric.
		argv := strings.Split(string(parts[1]), " ") // TODO: Handle argument quoting.
		if len(argv) < 1 {
			return nil, fmt.Errorf("process %q on line %d has invalid command", name, lineIndex)
		}
		process := Process{
			Name:      name,
			Command:   argv[0],
			Arguments: argv[1:],
		}
		procfile.Processes = append(procfile.Processes, process)
	}
	return &procfile, nil
}
