package inflect

import (
	"strings"
)

type Inflector struct {
	acronyms map[string]bool
}

func NewInflector() *Inflector {
	return &Inflector{
		acronyms: make(map[string]bool),
	}
}

var DefaultInflector = NewInflector()

func (infl *Inflector) AddAcronym(s string) {
	infl.acronyms[s] = true
}

func (infl *Inflector) IsAcronym(s string) bool {
	return infl.acronyms[s]
}

func init() {
	DefaultInflector.AddAcronym("id")
	DefaultInflector.AddAcronym("sid")
}

func KebabToPascal(s string) string {
	return KebabToGo(true, s)
}

func KebabToCamel(s string) string {
	return KebabToGo(false, s)
}

func KebabToGo(public bool, s string) string {
	return DefaultInflector.KebabToGo(public, s)
}

func (infl *Inflector) KebabToGo(public bool, s string) string {
	words := strings.Split(s, "-")
	up := public
	var b strings.Builder
	for _, word := range words {
		if up {
			if infl.IsAcronym(word) {
				_, _ = b.WriteString(strings.ToUpper(word))
			} else {
				_, _ = b.WriteString(strings.ToUpper(word[:1]))
				_, _ = b.WriteString(word[1:])
			}
		} else {
			_, _ = b.WriteString(word)
		}
		up = true
	}
	return b.String()
}
