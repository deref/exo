package manifest

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Names are restricted to be valid domain names.
func IsValidName(name string) bool {
	return ValidateName(name) == nil
}

func ValidateName(name string) error {
	if !namePattern.MatchString(name) {
		return fmt.Errorf("must match %q", namePattern)
	}
	if strings.Contains(name, "--") {
		return errors.New("cannot have two dashes in a row")
	}
	return nil
}

var namePattern = regexp.MustCompile("^[a-z]([-a-z0-9]*?[a-z0-9])?$")

func MangleName(name string) string {
	name = strings.ReplaceAll(name, "_", "-")
	name = invalidRegexp.ReplaceAllString(name, "")
	name = dashesRegexp.ReplaceAllString(name, "-")
	name = strings.TrimFunc(name, isInvalidChar)
	name = strings.TrimRight(name, "-")
	return name
}

var invalidRegexp = regexp.MustCompile("[^a-zA-Z0-9-]")

var dashesRegexp = regexp.MustCompile("-+")

func isInvalidChar(r rune) bool {
	return invalidRegexp.MatchString(string(r))
}
