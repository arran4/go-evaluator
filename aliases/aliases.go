package aliases

import "github.com/arran4/go-evaluator"

// Q is a short alias for evaluator.Query.
type Q = evaluator.Query

// Expr is a short alias for evaluator.Expression.
type Expr = evaluator.Expression

// EQ is a short alias for evaluator.IsExpression.
type EQ = evaluator.IsExpression

// NE is a short alias for evaluator.IsNotExpression.
type NE = evaluator.IsNotExpression

// GT is a short alias for evaluator.GreaterThanExpression.
type GT = evaluator.GreaterThanExpression

// GTE is a short alias for evaluator.GreaterThanOrEqualExpression.
type GTE = evaluator.GreaterThanOrEqualExpression

// LT is a short alias for evaluator.LessThanExpression.
type LT = evaluator.LessThanExpression

// LTE is a short alias for evaluator.LessThanOrEqualExpression.
type LTE = evaluator.LessThanOrEqualExpression

// CT is a short alias for evaluator.ContainsExpression.
type CT = evaluator.ContainsExpression

// AND is a short alias for evaluator.AndExpression.
type AND = evaluator.AndExpression

// OR is a short alias for evaluator.OrExpression.
type OR = evaluator.OrExpression

// NOT is a short alias for evaluator.NotExpression.
type NOT = evaluator.NotExpression
