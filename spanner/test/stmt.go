package test

import (
	"strings"

	"github.com/samber/lo"
)

func Split(stmts string) []string {
	result := []string{}
	for _, stmt := range strings.Split(stmts, ";") {
		lines := strings.Split(stmt, "\n")
		lines = lo.Filter(lines, func(line string, i int) bool {
			return !(strings.HasPrefix(line, "--") || strings.HasPrefix(line, "//"))
		})
		stmt = strings.Join(lines, " ")
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}
	return result
}
