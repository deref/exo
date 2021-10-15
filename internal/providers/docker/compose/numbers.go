package compose

import "strconv"

type Int String

func MakeInt(i int) Int {
	s := strconv.Itoa(i)
	return Int(String{
		Tag:        "!!int",
		Expression: s,
		Value:      s,
	})
}
