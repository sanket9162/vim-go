package highlight

import "regexp"

type TokenType int

const (
	TextToken TokenType = iota
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
