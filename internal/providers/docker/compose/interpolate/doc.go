// Package inteprolate implements recursive string substitution templating of
// yaml documents. The string template language is in the style of Python's
// string.Template library as extended by Docker Compose.
//
// References:
// - https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/interpolation.py
// - https://docs.python.org/3/library/string.html#string.Template
// - https://github.com/compose-spec/compose-spec/blob/b369fe5e02d80b619d14974cd1e64e7eea1b2345/spec.md#interpolation
package interpolate
