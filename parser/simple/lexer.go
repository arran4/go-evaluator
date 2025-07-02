package simple

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenIdent
	tokenString
	tokenNumber
	tokenAnd
	tokenOr
	tokenNot
	tokenIs
	tokenIsNot
	tokenContains
	tokenGT
	tokenGTE
	tokenLT
	tokenLTE
	tokenLParen
	tokenRParen
)

type token struct {
	typ tokenType
	val string
}

func isDelim(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
}

func lex(input string) ([]token, error) {
	var tokens []token
	i := 0
	for i < len(input) {
		r := rune(input[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}

		remain := input[i:]
		switch {
		case strings.HasPrefix(remain, "and") && (len(remain) == 3 || isDelim(rune(remain[3]))):
			tokens = append(tokens, token{typ: tokenAnd, val: "and"})
			i += 3
			continue
		case strings.HasPrefix(remain, "or") && (len(remain) == 2 || isDelim(rune(remain[2]))):
			tokens = append(tokens, token{typ: tokenOr, val: "or"})
			i += 2
			continue
		case strings.HasPrefix(remain, "not") && (len(remain) == 3 || isDelim(rune(remain[3]))):
			tokens = append(tokens, token{typ: tokenNot, val: "not"})
			i += 3
			continue
		case strings.HasPrefix(remain, "is not") && (len(remain) == 6 || isDelim(rune(remain[6]))):
			tokens = append(tokens, token{typ: tokenIsNot, val: "is not"})
			i += 6
			continue
		case strings.HasPrefix(remain, "is") && (len(remain) == 2 || isDelim(rune(remain[2]))):
			tokens = append(tokens, token{typ: tokenIs, val: "is"})
			i += 2
			continue
		case strings.HasPrefix(remain, "contains") && (len(remain) == 8 || isDelim(rune(remain[8]))):
			tokens = append(tokens, token{typ: tokenContains, val: "contains"})
			i += 8
			continue
		case strings.HasPrefix(remain, ">="):
			tokens = append(tokens, token{typ: tokenGTE, val: ">="})
			i += 2
			continue
		case strings.HasPrefix(remain, "<="):
			tokens = append(tokens, token{typ: tokenLTE, val: "<="})
			i += 2
			continue
		case strings.HasPrefix(remain, ">"):
			tokens = append(tokens, token{typ: tokenGT, val: ">"})
			i++
			continue
		case strings.HasPrefix(remain, "<"):
			tokens = append(tokens, token{typ: tokenLT, val: "<"})
			i++
			continue
		case strings.HasPrefix(remain, "("):
			tokens = append(tokens, token{typ: tokenLParen, val: "("})
			i++
			continue
		case strings.HasPrefix(remain, ")"):
			tokens = append(tokens, token{typ: tokenRParen, val: ")"})
			i++
			continue
		case remain[0] == '"':
			j := 1
			for i+j < len(input) && input[i+j] != '"' {
				j++
			}
			if i+j >= len(input) {
				return nil, fmt.Errorf("unterminated string")
			}
			tokens = append(tokens, token{typ: tokenString, val: input[i+1 : i+j]})
			i += j + 1
			continue
		default:
			if unicode.IsDigit(r) || (r == '.' && i+1 < len(input) && unicode.IsDigit(rune(input[i+1]))) {
				j := 1
				for i+j < len(input) && (unicode.IsDigit(rune(input[i+j])) || input[i+j] == '.') {
					j++
				}
				tokens = append(tokens, token{typ: tokenIdent, val: input[i : i+j]})
				i += j
				continue
			}
			j := 0
			for i+j < len(input) && !unicode.IsSpace(rune(input[i+j])) && !isDelim(rune(input[i+j])) {
				j++
			}
			if j == 0 {
				return nil, fmt.Errorf("unexpected character %q", input[i])
			}
			tokens = append(tokens, token{typ: tokenIdent, val: input[i : i+j]})
			i += j
			continue
		}
	}
	tokens = append(tokens, token{typ: tokenEOF})
	return tokens, nil
}
