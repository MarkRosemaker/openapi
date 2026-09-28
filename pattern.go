package openapi

import (
	"encoding/json/jsontext"
	"regexp"
	"strings"
)

// patternMarshal writes re in ECMA-262 syntax, the dialect of a Schema's "pattern".
func patternMarshal(enc *jsontext.Encoder, re *regexp.Regexp) error {
	return enc.WriteToken(jsontext.String(re2ToECMA(re.String())))
}

// patternUnmarshal compiles an ECMA-262 pattern with Go's regexp (RE2), translating the escapes RE2 spells differently.
func patternUnmarshal(dec *jsontext.Decoder, re *regexp.Regexp) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	compiled, err := regexp.Compile(ecmaToRE2(tkn.String()))
	if err != nil {
		return err
	}

	*re = *compiled

	return nil
}

// ecmaToRE2 rewrites ECMA-262's \uXXXX escapes as RE2's \x{XXXX}.
func ecmaToRE2(p string) string {
	return rewriteEscapes(p, func(rest string) (string, int) {
		if len(rest) >= 5 && rest[0] == 'u' && isHex(rest[1:5]) {
			return `\x{` + rest[1:5] + `}`, 5
		}

		return "", 0
	})
}

// re2ToECMA rewrites RE2's \x{X} to \x{XXXX} as ECMA-262's \uXXXX.
func re2ToECMA(p string) string {
	return rewriteEscapes(p, func(rest string) (string, int) {
		if len(rest) < 3 || rest[0] != 'x' || rest[1] != '{' {
			return "", 0
		}

		end := strings.IndexByte(rest, '}')
		if end < 3 || end > 6 || !isHex(rest[2:end]) {
			return "", 0
		}

		return `\u` + strings.Repeat("0", 6-end) + rest[2:end], end + 1
	})
}

// rewriteEscapes copies p, replacing an escape and the n bytes after its backslash where rewrite returns n > 0; an escaped backslash is never rewritten.
func rewriteEscapes(p string, rewrite func(rest string) (string, int)) string {
	var b strings.Builder

	for i := 0; i < len(p); i++ {
		if p[i] != '\\' || i+1 == len(p) {
			b.WriteByte(p[i])
			continue
		}

		if repl, n := rewrite(p[i+1:]); n > 0 {
			b.WriteString(repl)

			i += n

			continue
		}

		b.WriteString(p[i : i+2])
		i++
	}

	return b.String()
}

func isHex(s string) bool {
	for _, c := range s {
		if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F') {
			return false
		}
	}

	return s != ""
}
