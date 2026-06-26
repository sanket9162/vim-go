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

type RuleStruct struct {
	Regex *regexp.Regexp
	Type  TokenType
}
