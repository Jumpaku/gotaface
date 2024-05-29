package test

import (
	"github.com/Jumpaku/sqanner/tokenize"
	"github.com/samber/lo"
)

func Split(stmts string) []string {
	tokens, err := tokenize.Tokenize([]rune(stmts))
	if err != nil {
		panic(err)
	}
	result := []string{""}
	for _, token := range tokens {
		switch token.Kind {
		case tokenize.TokenComment:
			if result[len(result)-1] != "" {
				result = append(result, " ")
			}
			continue
		case tokenize.TokenSpecialChar:
			if string(token.Content) == ";" {
				result = append(result, "")
				continue
			}
		case tokenize.TokenSpace:
			if result[len(result)-1] == "" {
				continue
			}
		}
		result[len(result)-1] += string(token.Content)
	}

	return lo.Filter(result, func(item string, _ int) bool { return item != "" })
}
