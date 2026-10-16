package jsonutil

import "encoding/json/jsontext"

func nameSet(names []string) map[string]bool {
	set := make(map[string]bool, len(names))
	for _, name := range names {
		set[name] = true
	}

	return set
}

// memberName is the name of the object member p ends in, past the indices of any arrays it holds.
func memberName(p jsontext.Pointer) string {
	var name string

	for tok := range p.Tokens() {
		if !isIndex(tok) {
			name = tok
		}
	}

	return name
}

func isIndex(tok string) bool {
	if tok == "" {
		return false
	}

	for _, r := range tok {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
