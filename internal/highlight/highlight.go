package highlight

import (
	"regexp"
)

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
	Group int // The capturing group to highlight (0 means the whole match)
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
	{Regex: regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`), Type: TokenFunction, Group: 1},
}

// TokenizeLine returns a slice mapping each rune index of the line to its TokenType
func TokenizeLine(line string, rules []Rule) []TokenType {
	runes := []rune(line)
	tokens := make([]TokenType, len(runes))
	for i := range tokens {
		tokens[i] = TokenText
	}

	if len(line) == 0 {
		return tokens
	}

	// Map byte index to rune index for UTF-8 safety
	byteToRune := make([]int, len(line)+1)
	runeIdx := 0
	for byteIdx, r := range line {
		byteToRune[byteIdx] = runeIdx
		runeLen := len(string(r))
		for i := 1; i < runeLen; i++ {
			byteToRune[byteIdx+i] = runeIdx
		}
		runeIdx++
	}
	byteToRune[len(line)] = runeIdx

	// Apply rules in reverse order of precedence (or handle priority)
	for _, rule := range rules {
		matches := rule.Regex.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			startByte, endByte := match[0], match[1]

			// Highlight the specified capture group if specified
			if rule.Group > 0 && len(match) >= 2*(rule.Group+1) {
				gStart := match[2*rule.Group]
				gEnd := match[2*rule.Group+1]
				if gStart != -1 && gEnd != -1 {
					startByte, endByte = gStart, gEnd
				}
			}

			startRune := byteToRune[startByte]
			endRune := byteToRune[endByte]

			for idx := startRune; idx < endRune; idx++ {
				tokens[idx] = rule.Type
			}
		}
	}
	return tokens
}
