package openapi

import (
	"fmt"
	"net/http"
	"strconv"
)

// StatusCode is used as a key in the Responses map to describe the expected response for that HTTP status code.
// Any HTTP status code can be used, but only one property per code, to describe the expected response for that HTTP status code.
// To define a range of response codes, this field MAY contain the uppercase wildcard character `X`. For example, `2XX` represents all response codes between `[200-299]`. Only the following range definitions are allowed: `1XX`, `2XX`, `3XX`, `4XX`, and `5XX`. If a response is defined using an explicit code, the explicit code definition takes precedence over the range definition for that code.
type StatusCode string

const (
	// Use default to documentat responses other than the ones declared for specific HTTP response codes. Use as a field in the Responses map to cover undeclared responses.
	StatusCodeDefault StatusCode = "default"
)

// StatusText returns the text representation of the status code.
// If the status code is `"default"` or an invalid status code, it returns an empty string.
func (sc StatusCode) StatusText() string {
	switch sc {
	case "default":
		return ""
	default:
		code, _ := strconv.Atoi(string(sc))
		return http.StatusText(code)
	}
}

func (sc StatusCode) Validate() error {
	if sc == StatusCodeDefault {
		return nil
	}

	// Check if the status code is a range definition,
	// e.g. `1XX`, `2XX`, `3XX`, `4XX`, and `5XX`.
	if len(sc) == 3 && sc[1] == 'X' && sc[2] == 'X' {
		switch sc[0] {
		case '1', '2', '3', '4', '5':
			return nil
		}
	}

	code, err := strconv.Atoi(string(sc))
	if err != nil || http.StatusText(code) == "" {
		return fmt.Errorf("invalid status code %q", sc)
	}

	return nil
}

func (sc StatusCode) IsSuccess() bool {
	// Check if the status code is a range definition,
	// e.g. `1XX`, `2XX`, `3XX`, `4XX`, and `5XX`.
	if len(sc) == 3 && sc[1] == 'X' && sc[2] == 'X' {
		return sc[0] == '2'
	}

	code, _ := strconv.Atoi(string(sc))
	return 200 <= code && code < 300
}
