package openapi

import (
	"fmt"
	"regexp"
)

var reTemplateExpressions = regexp.MustCompile(`\{([^}]+)\}`)

// ParsedPath is a parsed URL path template.
type ParsedPath struct {
	// Format string with the variables replaced with %s.
	Format string
	// VariableNames is the list of variable keys in order.
	VariableNames []string
}

func (path Path) Parse() ParsedPath {
	p := ParsedPath{}

	// Find all variables in the url
	p.Format = reTemplateExpressions.ReplaceAllStringFunc(string(path), func(s string) string {
		// store the variable name
		p.VariableNames = append(p.VariableNames, s[1:len(s)-1])

		return "%s"
	})

	return p
}

func (p ParsedPath) String() string {
	vars := make([]any, len(p.VariableNames))
	for i, v := range p.VariableNames {
		vars[i] = fmt.Sprintf("{%s}", v)
	}

	return fmt.Sprintf(p.Format, vars...)
}
