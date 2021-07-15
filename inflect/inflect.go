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
	return infl.acronyms[infl.Singularize(s)]
}

func init() {
	DefaultInflector.AddAcronym("id")
	DefaultInflector.AddAcronym("sid")
}

func KebabToPublic(s string) string {
	return KebabToGo(true, s)
}

func KebabToPrivate(s string) string {
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
			word = infl.LowerToTitle(word)
		}
		_, _ = b.WriteString(word)
		up = true
	}
	return b.String()
}

func KebabToJSVar(s string) string {
	return DefaultInflector.KebabToJS(false, s)
}

func (infl *Inflector) KebabToJS(class bool, s string) string {
	words := strings.Split(s, "-")
	up := class
	var b strings.Builder
	for _, word := range words {
		if up {
			word = strings.ToUpper(word[:1]) + word[1:]
		}
		_, _ = b.WriteString(word)
		up = true
	}
	return b.String()
}

func (infl *Inflector) LowerToTitle(word string) string {
	if infl.IsPluralAcronym(word) {
		return strings.ToUpper(word[:len(word)-1]) + "s" // XXX
	} else if infl.IsAcronym(word) {
		return strings.ToUpper(word)
	} else {
		return strings.ToUpper(word[:1]) + word[1:]
	}
}

func (infl *Inflector) IsPluralAcronym(word string) bool {
	if infl.IsSingular(word) {
		return false
	}
	return infl.IsAcronym(infl.Singularize(word))
}

func (infl *Inflector) Singularize(word string) string {
	if strings.HasSuffix(word, "s") {
		return word[0 : len(word)-1]
	}
	return word
}

func (infl *Inflector) IsSingular(word string) bool {
	return !infl.IsPlural(word)
}

func (infl *Inflector) IsPlural(word string) bool {
	return strings.HasSuffix(word, "s") // XXX
}
