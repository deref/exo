package manifest

import "regexp"

var NamePattern = regexp.MustCompile("^[a-z]([-a-z0-9]*?[a-z0-9])?$")

func IsValidName(name string) bool {
	return NamePattern.MatchString(name)
}
