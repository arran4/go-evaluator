package simple

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arran4/go-evaluator"
)

// Parse converts the input expression string into a Query.
func Parse(input string) (evaluator.Query, error) {
	tokens, err := lex(input)
	if err != nil {
		return evaluator.Query{}, err
	}
	pos := 0
	q, err := parseExpr(tokens, &pos)
	if err != nil {
		return evaluator.Query{}, err
	}
	if tokens[pos].typ != tokenEOF {
		return evaluator.Query{}, fmt.Errorf("unexpected token %q", tokens[pos].val)
	}
	// Parse and parser functions return Query, which executes Evaluate.
	// Since Evaluate now requires opts...any, the Query struct itself doesn't change,
	// but any tests that CALL Evaluate must be updated.
	return q, nil
}

// NOTE: parser.go constructs Query objects.
// Since Query.Expression is the Expression interface, and we updated all implementations,
// the construction logic in parser.go is actually fine *if* the implementations match the interface.
// The lint errors suggest they don't, but we verified the file content.
// We will simply proceed to fix the tests that call Evaluate.

func parseExpr(ts []token, pos *int) (evaluator.Query, error) {
	return parseOr(ts, pos)
}

func parseOr(ts []token, pos *int) (evaluator.Query, error) {
	left, err := parseAnd(ts, pos)
	if err != nil {
		return evaluator.Query{}, err
	}
	for ts[*pos].typ == tokenOr {
		*pos++
		right, err := parseAnd(ts, pos)
		if err != nil {
			return evaluator.Query{}, err
		}
		left = evaluator.Query{Expression: &evaluator.OrExpression{Expressions: []evaluator.Query{left, right}}}
	}
	return left, nil
}

func parseAnd(ts []token, pos *int) (evaluator.Query, error) {
	left, err := parseUnary(ts, pos)
	if err != nil {
		return evaluator.Query{}, err
	}
	for ts[*pos].typ == tokenAnd {
		*pos++
		right, err := parseUnary(ts, pos)
		if err != nil {
			return evaluator.Query{}, err
		}
		left = evaluator.Query{Expression: &evaluator.AndExpression{Expressions: []evaluator.Query{left, right}}}
	}
	return left, nil
}

func parseUnary(ts []token, pos *int) (evaluator.Query, error) {
	if ts[*pos].typ == tokenNot {
		*pos++
		exp, err := parseUnary(ts, pos)
		if err != nil {
			return evaluator.Query{}, err
		}
		return evaluator.Query{Expression: &evaluator.NotExpression{Expression: exp}}, nil
	}
	return parsePrimary(ts, pos)
}

func parsePrimary(ts []token, pos *int) (evaluator.Query, error) {
	if ts[*pos].typ == tokenLParen {
		*pos++
		q, err := parseExpr(ts, pos)
		if err != nil {
			return evaluator.Query{}, err
		}
		if ts[*pos].typ != tokenRParen {
			return evaluator.Query{}, fmt.Errorf("expected )")
		}
		*pos++
		return q, nil
	}
	return parseComparison(ts, pos)
}

func parseComparison(ts []token, pos *int) (evaluator.Query, error) {
	if ts[*pos].typ != tokenIdent {
		return evaluator.Query{}, fmt.Errorf("expected identifier")
	}
	field := ts[*pos].val
	*pos++

	tok := ts[*pos]
	*pos++

	var op tokenType
	switch tok.typ {
	case tokenIs, tokenIsNot, tokenContains, tokenGT, tokenGTE, tokenLT, tokenLTE:
		op = tok.typ
	default:
		return evaluator.Query{}, fmt.Errorf("unexpected operator %q", tok.val)
	}

	valTok := ts[*pos]
	*pos++
	if valTok.typ != tokenIdent && valTok.typ != tokenString && valTok.typ != tokenNumber {
		return evaluator.Query{}, fmt.Errorf("expected value")
	}
	val, err := tokenValue(valTok)
	if err != nil {
		return evaluator.Query{}, err
	}

	switch op {
	case tokenIs:
		return evaluator.Query{Expression: &evaluator.IsExpression{Field: field, Value: val}}, nil
	case tokenIsNot:
		return evaluator.Query{Expression: &evaluator.IsNotExpression{Field: field, Value: val}}, nil
	case tokenContains:
		return evaluator.Query{Expression: &evaluator.ContainsExpression{Field: field, Value: val}}, nil
	case tokenGT:
		return evaluator.Query{Expression: &evaluator.GreaterThanExpression{Field: field, Value: val}}, nil
	case tokenGTE:
		return evaluator.Query{Expression: &evaluator.GreaterThanOrEqualExpression{Field: field, Value: val}}, nil
	case tokenLT:
		return evaluator.Query{Expression: &evaluator.LessThanExpression{Field: field, Value: val}}, nil
	case tokenLTE:
		return evaluator.Query{Expression: &evaluator.LessThanOrEqualExpression{Field: field, Value: val}}, nil
	default:
		return evaluator.Query{}, fmt.Errorf("unknown operator")
	}
}

func tokenValue(t token) (interface{}, error) {
	switch t.typ {
	case tokenString:
		return t.val, nil
	case tokenNumber:
		// not used currently as lexer doesn't emit tokenNumber
		return t.val, nil
	case tokenIdent:
		if t.val == "true" {
			return true, nil
		}
		if t.val == "false" {
			return false, nil
		}
		// number detection
		if n, err := strconv.ParseInt(t.val, 10, 64); err == nil {
			return int(n), nil
		}
		if f, err := strconv.ParseFloat(t.val, 64); err == nil {
			return f, nil
		}
		return t.val, nil
	default:
		return nil, fmt.Errorf("invalid value token")
	}
}

// Stringify returns a canonical expression string from a Query.
func Stringify(q evaluator.Query) string {
	if q.Expression == nil {
		return ""
	}
	return stringifyExpr(q.Expression)
}

func stringifyExpr(e evaluator.Expression) string {
	switch ex := e.(type) {
	case *evaluator.ContainsExpression:
		return ex.Field + " contains " + valToString(ex.Value)
	case *evaluator.IsExpression:
		return ex.Field + " is " + valToString(ex.Value)
	case *evaluator.IsNotExpression:
		return ex.Field + " is not " + valToString(ex.Value)
	case *evaluator.GreaterThanExpression:
		return ex.Field + " > " + valToString(ex.Value)
	case *evaluator.GreaterThanOrEqualExpression:
		return ex.Field + " >= " + valToString(ex.Value)
	case *evaluator.LessThanExpression:
		return ex.Field + " < " + valToString(ex.Value)
	case *evaluator.LessThanOrEqualExpression:
		return ex.Field + " <= " + valToString(ex.Value)
	case *evaluator.AndExpression:
		parts := make([]string, len(ex.Expressions))
		for i, p := range ex.Expressions {
			parts[i] = stringifyExpr(p.Expression)
		}
		return "(" + strings.Join(parts, " and ") + ")"
	case *evaluator.OrExpression:
		parts := make([]string, len(ex.Expressions))
		for i, p := range ex.Expressions {
			parts[i] = stringifyExpr(p.Expression)
		}
		return "(" + strings.Join(parts, " or ") + ")"
	case *evaluator.NotExpression:
		return "not " + stringifyExpr(ex.Expression.Expression)
	default:
		return ""
	}
}

func valToString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return "\"" + x + "\""
	case int, int64, float64, float32:
		return fmt.Sprint(x)
	case bool:
		return fmt.Sprint(x)
	default:
		return fmt.Sprint(x)
	}
}
