package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"
)

// bs is a backslash, so that \u escapes below are built rather than typed.
const bs = `\`

func TestPatternEscapes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ ecma, re2 string }{
		{`^[^` + bs + `u0000-` + bs + `u001f` + bs + `u007f]+$`, `^[^\x{0000}-\x{001f}\x{007f}]+$`},
		{bs + `uABCD`, `\x{ABCD}`},
		// an escaped backslash followed by "u" is a literal backslash and "u".
		{bs + bs + `u0041`, bs + bs + `u0041`},
		{bs + bs + bs + `u0041`, `\\\x{0041}`},
		// anything else is left as it is.
		{`^[a-z]+\d{2}$`, `^[a-z]+\d{2}$`},
		{bs + `u00`, bs + `u00`},
		{`trailing\`, `trailing\`},
	} {
		if got := ecmaToRE2(tc.ecma); got != tc.re2 {
			t.Errorf("ecmaToRE2(%q) = %q, want %q", tc.ecma, got, tc.re2)
		}

		if got := re2ToECMA(tc.re2); got != tc.ecma {
			t.Errorf("re2ToECMA(%q) = %q, want %q", tc.re2, got, tc.ecma)
		}
	}

	// RE2 escapes of up to 4 hex digits are padded; longer ones have no \u form.
	for re2, ecma := range map[string]string{
		`\x{41}`:     bs + `u0041`,
		`\x{10FFFF}`: `\x{10FFFF}`,
	} {
		if got := re2ToECMA(re2); got != ecma {
			t.Errorf("re2ToECMA(%q) = %q, want %q", re2, got, ecma)
		}
	}
}

// TestSchema_PatternJSON round-trips an ECMA-262 pattern that Go's regexp
// only accepts in another syntax: it compiles, and is written back as read.
func TestSchema_PatternJSON(t *testing.T) {
	t.Parallel()

	in := `{"type":"string","pattern":"^[^` + bs + bs + `u0000-` + bs + bs + `u001f` + bs + bs + `u007f]+$"}`

	var s Schema
	if err := json.Unmarshal([]byte(in), &s, jsonOpts); err != nil {
		t.Fatal(err)
	}

	if !s.Pattern.MatchString("page title") || s.Pattern.MatchString("line\nbreak") {
		t.Errorf("pattern %q matches the wrong strings", s.Pattern)
	}

	out, err := json.Marshal(&s, jsonOpts)
	if err != nil {
		t.Fatal(err)
	}

	got := jsontext.Value(out)
	if err := got.Compact(); err != nil {
		t.Fatal(err)
	}

	if string(got) != in {
		t.Errorf("got %s, want %s", got, in)
	}
}
