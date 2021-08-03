package strutil

import "bytes"

func EscapeString(in string, specialChars []rune, escapeChar rune) string {
	out := bytes.NewBuffer(make([]byte, 0, len(in)))
	for _, ch := range in {
		if ch == escapeChar {
			out.WriteRune(escapeChar)
		} else {
			for _, specialChar := range specialChars {
				if ch == specialChar {
					out.WriteRune(escapeChar)
					break
				}
			}
		}
		out.WriteRune(ch)
	}

	return out.String()
}

func UnescapeString(in string, escapeChar rune) string {
	out := bytes.NewBuffer(make([]byte, 0, len(in)))
	var isEscaped = false
	for _, ch := range in {
		if ch == escapeChar && !isEscaped {
			isEscaped = true
		} else {
			isEscaped = false
			out.WriteRune(ch)
		}
	}

	return out.String()
}
