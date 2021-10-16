// Package compose implements an AST, parser, validator, and interpolator for
// <https://github.com/compose-spec/compose-spec>.
//
// The AST preserves some formatting information from the source Yaml. This
// enables extracting subdocuments in a way that perserves source formatting,
// and in particular maintains the stable ordering of map entries.
//
// Parsing is implemented in terms of yaml unmarshaling and performs partial
// evaluation of interpolation. All scalar leaf nodes contain a String
// structure which maintains the original Expression string and an interpolated
// Value string. Other scalar types (numbers and booleans), as well as
// composite types, will embed a String to facitate delaying interpolation
// until an environment is available.
package compose
