package highlight

import "regexp"

type TokenType int

const (
	TokenText TokenType = iota
	TokenKeyword
	TokenTypeWord
	TokenString
	TokenComment
	TokenNumber
	TokenFunction
	TokenOperator
)

type Rule struct {
	Regex *regexp.Regexp
	Type  TokenType
}

// Highlighting rules for Go syntax
var GoRules = []Rule{
	{Regex: regexp.MustCompile(`//.*`), Type: TokenComment},
	{Regex: regexp.MustCompile(`"([^"\\]|\\.)*"`), Type: TokenString},
	{Regex: regexp.MustCompile("`[^`]*`"), Type: TokenString},
	{Regex: regexp.MustCompile(`'([^'\\]|\\.)*'`), Type: TokenString},
	{Regex: regexp.MustCompile(`\b[0-9]+(\.[0-9]+)?\b`), Type: TokenNumber},
	{Regex: regexp.MustCompile(`\b(func|package|import|return|if|else|for|range|switch|case|default|type|struct|interface|go|select|chan|map|var|const|defer)\b`), Type: TokenKeyword},
	{Regex: regexp.MustCompile(`\b(int|string|bool|float64|float32|rune|byte|error|uint|uintptr|nil|true|false)\b`), Type: TokenTypeWord},
	{Regex: regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\s*(?=\()`), Type: TokenFunction},
}

// TokenizeLine returns a slice mapping each rune index of the line to its TokenType
func TokenizeLine(line string, rules []Rule) []TokenType {
	tokens := make([]TokenType, len(line))
	for i := range tokens {
		tokens[i] = TokenText
	}

	// Apply rules in reverse order of precedence (or handle priority)
	for _, rule := range rules {
		matches := rule.Regex.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			start, end := match[0], match[1]
			for idx := start; idx < end; idx++ {
				tokens[idx] = rule.Type
			}
		}
	}
	return tokens
}
