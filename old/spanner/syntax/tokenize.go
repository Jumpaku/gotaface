package syntax

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
	"unicode"
)

type TokenCode int

const (
	TokenUnknown TokenCode = iota
	TokenEOF
	TokenSpace
	TokenComment
	TokenIdentifier
	TokenIdentifierQuoted
	TokenLiteralQuoted
	TokenLiteralInteger
	TokenLiteralFloat
	TokenKeyword
	TokenSpecialChar
)

const specialChars = "@,()[]{}<>.;:/+-~*/|&^=!"
const additionalSpecialChars = "$?"

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
func isLetter(r rune) bool {
	return r == '_' || 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
}
func isSpecial(r rune) bool {
	return strings.ContainsRune(specialChars, r)
}

var keywords = map[string]bool{
	"ALL":                  true,
	"AND":                  true,
	"ANY":                  true,
	"ARRAY":                true,
	"AS":                   true,
	"ASC":                  true,
	"ASSERT_ROWS_MODIFIED": true,
	"AT":                   true,
	"BETWEEN":              true,
	"BY":                   true,
	"CASE":                 true,
	"CAST":                 true,
	"COLLATE":              true,
	"CONTAINS":             true,
	"CREATE":               true,
	"CROSS":                true,
	"CUBE":                 true,
	"CURRENT":              true,
	"DEFAULT":              true,
	"DEFINE":               true,
	"DESC":                 true,
	"DISTINCT":             true,
	"ELSE":                 true,
	"END":                  true,
	"ENUM":                 true,
	"ESCAPE":               true,
	"EXCEPT":               true,
	"EXCLUDE":              true,
	"EXISTS":               true,
	"EXTRACT":              true,
	"FALSE":                true,
	"FETCH":                true,
	"FOLLOWING":            true,
	"FOR":                  true,
	"FROM":                 true,
	"FULL":                 true,
	"GROUP":                true,
	"GROUPING":             true,
	"GROUPS":               true,
	"HASH":                 true,
	"HAVING":               true,
	"IF":                   true,
	"IGNORE":               true,
	"IN":                   true,
	"INNER":                true,
	"INTERSECT":            true,
	"INTERVAL":             true,
	"INTO":                 true,
	"IS":                   true,
	"JOIN":                 true,
	"LATERAL":              true,
	"LEFT":                 true,
	"LIKE":                 true,
	"LIMIT":                true,
	"LOOKUP":               true,
	"MERGE":                true,
	"NATURAL":              true,
	"NEW":                  true,
	"NO":                   true,
	"NOT":                  true,
	"NULL":                 true,
	"NULLS":                true,
	"OF":                   true,
	"ON":                   true,
	"OR":                   true,
	"ORDER":                true,
	"OUTER":                true,
	"OVER":                 true,
	"PARTITION":            true,
	"PRECEDING":            true,
	"PROTO":                true,
	"RANGE":                true,
	"RECURSIVE":            true,
	"RESPECT":              true,
	"RIGHT":                true,
	"ROLLUP":               true,
	"ROWS":                 true,
	"SELECT":               true,
	"SET":                  true,
	"SOME":                 true,
	"STRUCT":               true,
	"TABLESAMPLE":          true,
	"THEN":                 true,
	"TO":                   true,
	"TREAT":                true,
	"TRUE":                 true,
	"UNBOUNDED":            true,
	"UNION":                true,
	"UNNEST":               true,
	"USING":                true,
	"WHEN":                 true,
	"WHERE":                true,
	"WINDOW":               true,
	"WITH":                 true,
	"WITHIN":               true,

	"NUMERIC":   true,
	"DATE":      true,
	"TIMESTAMP": true,
	"JSON":      true,
}

type Token struct {
	Code    TokenCode
	Content string

	Begin  int
	End    int
	Line   int
	Column int
}

type TokenScanner struct {
	Input   []rune
	Cursor  int
	lines   int
	columns int
}

func (s *TokenScanner) ScanNext() (Token, error) {
	return Token{}, nil
}

func (s *TokenScanner) accept(n int, code TokenCode) Token {
	out := s.Input[s.Cursor : s.Cursor+n]
	token := Token{
		Code:    code,
		Content: string(out),
		Begin:   s.Cursor,
		End:     s.Cursor + n,
		Line:    s.lines,
		Column:  s.columns,
	}
	for _, o := range out {
		if o == '\n' {
			s.lines++
			s.columns = 0
		}
		s.columns++
		s.Cursor++
	}
	return token
}

func (s *TokenScanner) len() int {
	return len(s.Input) - s.Cursor
}
func (s *TokenScanner) peekAt(n int) rune {
	return s.Input[s.Cursor+n]
}
func (s *TokenScanner) peekSlice(begin int, endExclusive int) []rune {
	return s.Input[s.Cursor+begin : s.Cursor+endExclusive]
}
func (s *TokenScanner) countWhile(begin int, satisfy func(rune) bool) int {
	count := 0
	for cur := s.Cursor + begin; cur < len(s.Input); cur++ {
		if satisfy(s.Input[cur]) {
			count++
		} else {
			break
		}
	}
	return count
}
func (s *TokenScanner) findFirst(begin int, patternSize int, pattern func([]rune) bool) (int, bool) {
	for cur := s.Cursor + begin; cur < len(s.Input)-patternSize; cur++ {
		if pattern(s.Input[cur : cur+patternSize]) {
			return cur - s.Cursor, true
		}
	}

	return len(s.Input[s.Cursor:]) - s.Cursor, false
}

func (s *TokenScanner) wrapErr(err error) error {
	sizeAfter := 20
	if sizeAfter > s.len() {
		sizeAfter = s.len()
	}
	sizeBefore := 5
	if sizeBefore > s.Cursor {
		sizeBefore = s.Cursor
	}
	input := string(s.Input[s.Cursor-sizeBefore : s.Cursor+sizeAfter])
	return fmt.Errorf(`fail to scan token at line %d column %d near ...%s...: %w`, s.lines, s.columns, input, err)
}

type nextFunc func(s *TokenScanner) (f nextFunc, done bool, err error)

func Tokenize(input []rune) ([]Token, error) {
	s := &TokenScanner{Input: input}

	next := initial
	for {
		var done bool
		var err error
		next, done, err = next(s)
		if err != nil {
			return nil, fmt.Errorf(`fail to read token: %w`, err)
		}
		if done {
			return s.tokens, nil
		}
	}
}

func initial(s *TokenScanner) (nextFunc, bool, error) {
	if s.Cursor >= len(s.Input) {
		s.consumeToken(0, TokenEOF)
		return nil, true, nil
	}
	switch {
	case unicode.IsSpace(s.peekAt(0)):
		return spaces, false, nil
	case s.peekAt(0) == '#',
		slices.Contains([]string{`/*`, `//`, `--`}, string(s.peekSlice(0, 2))):
		return comment, false, nil
	case s.peekAt(0) == '`':
		return identifierQuoted, false, nil
	case s.peekAt(0) == '"', s.peekAt(0) == '\'',
		slices.Contains([]string{`r"`, `r'`, `b"`, `b'`}, strings.ToLower(string(s.peekSlice(0, 2)))),
		slices.Contains([]string{`rb"`, `rb'`, `br"`, `br'`}, strings.ToLower(string(s.peekSlice(0, 2)))):
		return identifierQuoted, false, nil
	case isDigit(s.peekAt(0)):
		n := s.countWhile(0, unicode.IsDigit)
		if s.peekAt(n) == '.' {
		}
		m := s.countIf(func(r rune) bool { return strings.ContainsRune(digitChars, r) })
		if s.peekAt(n) == '.' {
		}
	case s.peekAt(0) == '.':
		if strings.ContainsRune(digitChars, s.peekAt(1)) {
		}
		return specialChar, false, nil
	case strings.ContainsRune(specialChars, s.peekAt(0)):
		return specialChar, false, nil
	case strings.ContainsRune(letterChars, s.peekAt(0)):
		n := s.countIf(func(r rune) bool { return strings.ContainsRune(letterChars+digitChars, r) })
		lowerSlice := strings.ToUpper(string(s.peekSlice(0, n)))
		if keywords[lowerSlice] {
			return keyword, false, nil
		}
		return identifier, false, nil
	}
	return spaces, false, nil
}

func spaces(s *TokenScanner) (nextFunc, bool, error) {
	n := s.countIf(unicode.IsSpace)
	s.consumeToken(n, TokenSpace)
	return initial, false, nil
}
func comment(s *TokenScanner) (nextFunc, bool, error) {
	switch {
	default:
		return nil, false, fmt.Errorf(`fail to find beginning of comment: '/*', '//', '--', or '#' is expected`)
	case s.hasPrefix(`/*`, true):
		n, ok := s.findFirst(2, 2, func(s string) bool { return s == `*/` })
		if !ok {
			return nil, false, fmt.Errorf(`fail to find end of comment: '*/' is expected`)
		}

		s.consumeToken(n, TokenComment)
		return initial, false, nil
	case s.hasPrefix(`//`, true), s.hasPrefix(`--`, true), s.hasPrefix(`#`, true):
		n, _ := s.findFirst(1, 1, func(s string) bool { return s == "\n" })
		s.consumeToken(n, TokenComment)
		return initial, false, nil
	}
}
func specialChar(s *TokenScanner) (nextFunc, bool, error) {
	s.consumeToken(1, TokenSpecialChar)
	return initial, false, nil
}

func literalQuoted(s *TokenScanner) (nextFunc, bool, error) {
	n, ok := s.findFirst(1, 1, func(s string) bool { return s == "`" })
	if !ok {
		return nil, false, fmt.Errorf("fail to find end of back quote: '`' is expected")
	}

	s.consumeToken(n, TokenComment)
	return initial, false, nil
}

func keyword(s *TokenScanner) (nextFunc, bool, error) {
	slices.Index(lowerKeywords)
	s.consumeToken(n, TokenSpace)
	return initial, false, nil
}

func identifierQuoted(s *TokenScanner) (nextFunc, bool, error) {
	n, ok := s.findFirst(1, 1, func(s string) bool { return s == "`" })
	if !ok {
		return nil, false, fmt.Errorf("fail to find end of back quote: '`' is expected")
	}

	s.consumeToken(n, TokenComment)
	return initial, false, nil
}

func identifier(s *TokenScanner) (nextFunc, bool, error) {
	s.peekAt(0)
	s.consumeToken(1, TokenSpecialChar)
	return initial, false, nil
}
