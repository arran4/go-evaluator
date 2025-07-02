// Package evaluator provides a simple expression language for querying Go
// structs. Expressions are represented as Go structs which can be evaluated
// against arbitrary values or composed together using logical operators. The
// package also supports marshaling and unmarshaling expressions to and from
// JSON for storage or transmission.
package evaluator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// derefValue dereferences pointer inputs and returns the underlying value.
// It supports structs and maps and returns false for all other types.
func derefValue(i interface{}) (reflect.Value, bool) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}, false
		}
		v = v.Elem()
	} else if v.Kind() == reflect.Struct {
		// maintain backward behaviour: require pointer for structs
		return reflect.Value{}, false
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Map:
		return v, true
	default:
		return reflect.Value{}, false
	}
}

// getField retrieves a field value from either a struct or map value.
// For structs it uses FieldByName, while for maps it looks up the key by name.
func getField(v reflect.Value, name string) (reflect.Value, bool) {
	switch v.Kind() {
	case reflect.Struct:
		f := v.FieldByName(name)
		if f.IsValid() {
			return f, true
		}
		return reflect.Value{}, false
	case reflect.Map:
		key := reflect.ValueOf(name)
		if key.Type().AssignableTo(v.Type().Key()) {
			f := v.MapIndex(key)
			if f.IsValid() {
				return f, true
			}
		}
		return reflect.Value{}, false
	default:
		return reflect.Value{}, false
	}
}

// Expression represents a single boolean expression that can be evaluated
// against a struct value.
type Expression interface {
	// Evaluate returns true if the expression matches the supplied value.
	Evaluate(i interface{}) bool
}

// ContainsExpression checks whether a slice field contains the given Value.
type ContainsExpression struct {
	Field string
	Value interface{}
}

func (e ContainsExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	if f.Type().Kind() != reflect.Slice {
		return false
	}
	cv := reflect.ValueOf(e.Value)
	if !cv.IsValid() {
		return false
	}
	if f.Type().Elem().Kind() != cv.Type().Kind() {
		return false
	}
	for i := 0; i < f.Len(); i++ {
		if reflect.DeepEqual(f.Index(i).Interface(), cv.Interface()) {
			return true
		}
	}
	return false
}

// IsNotExpression succeeds when the specified Field does not equal Value.
type IsNotExpression struct {
	Field string
	Value interface{}
}

func (e IsNotExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	return !reflect.DeepEqual(f.Interface(), e.Value)
}

// IsExpression succeeds when the specified Field equals Value.
type IsExpression struct {
	Field string
	Value interface{}
}

func (e IsExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	if e.Value == nil {
		switch f.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
			if f.IsNil() {
				return true
			}
		}
	}
	return reflect.DeepEqual(f.Interface(), e.Value)
}

// AndExpression evaluates to true only if all child Expressions do as well.
type AndExpression struct {
	Expressions []Query `json:"Expressions"`
}

func (e AndExpression) Evaluate(i interface{}) bool {
	for _, q := range e.Expressions {
		if !q.Evaluate(i) {
			return false
		}
	}
	return true
}

// OrExpression evaluates to true if any of the child Expressions do.
type OrExpression struct {
	Expressions []Query `json:"Expressions"`
}

func (e OrExpression) Evaluate(i interface{}) bool {
	for _, q := range e.Expressions {
		if q.Evaluate(i) {
			return true
		}
	}
	return false
}

// NotExpression inverts the result of a single child Expression.
type NotExpression struct {
	Expression Query `json:"Expression"`
}

func (e NotExpression) Evaluate(i interface{}) bool {
	return !e.Expression.Evaluate(i)
}

func numericValue(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case uintptr:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		if err == nil {
			return f, true
		}
		return 0, false
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func stringValue(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	default:
		return fmt.Sprint(v)
	}
}

// GreaterThanExpression compares Field to Value and succeeds when the field is
// greater than the provided value.
type GreaterThanExpression struct {
	Field string
	Value interface{}
}

func (e GreaterThanExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Int()) > v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Uint()) > v
	case reflect.Float32, reflect.Float64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return f.Float() > v
	case reflect.String:
		sval := stringValue(e.Value)
		return strings.Compare(f.String(), sval) > 0
	default:
		return false
	}
}

// GreaterThanOrEqualExpression succeeds when Field is greater than or equal to
// Value.
type GreaterThanOrEqualExpression struct {
	Field string
	Value interface{}
}

func (e GreaterThanOrEqualExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Int()) >= v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Uint()) >= v
	case reflect.Float32, reflect.Float64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return f.Float() >= v
	case reflect.String:
		sval := stringValue(e.Value)
		return strings.Compare(f.String(), sval) >= 0
	default:
		return false
	}
}

// LessThanExpression succeeds when Field is strictly less than Value.
type LessThanExpression struct {
	Field string
	Value interface{}
}

func (e LessThanExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Int()) < v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Uint()) < v
	case reflect.Float32, reflect.Float64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return f.Float() < v
	case reflect.String:
		sval := stringValue(e.Value)
		return strings.Compare(f.String(), sval) < 0
	default:
		return false
	}
}

// LessThanOrEqualExpression succeeds when Field is less than or equal to Value.
type LessThanOrEqualExpression struct {
	Field string
	Value interface{}
}

func (e LessThanOrEqualExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Int()) <= v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return float64(f.Uint()) <= v
	case reflect.Float32, reflect.Float64:
		v, ok := numericValue(e.Value)
		if !ok {
			return false
		}
		return f.Float() <= v
	case reflect.String:
		sval := stringValue(e.Value)
		return strings.Compare(f.String(), sval) <= 0
	default:
		return false
	}
}

// QueryRaw is the JSON representation of a query. ExpressionRawJson stores the
// raw JSON for the underlying expression and is resolved during unmarshalling.
type QueryRaw struct {
	Expression        Expression      `json:"-"`
	ExpressionRawJson json.RawMessage `json:"Expression"`
}

// Query wraps QueryRaw and provides evaluation and JSON unmarshalling helpers.
type Query QueryRaw

// typedExpression couples an Expression value with a Type field
// so it can be marshaled and unmarshaled in a generic fashion.
// The Expression field is strongly typed using generics.
type typedExpression[E Expression] struct {
	Type       string `json:"Type"`
	Expression E      `json:"Expression"`
}

// marshalExpression serializes any Expression along with its type
// indicator using typedExpression.
func marshalExpression(e Expression) ([]byte, error) {
	switch expr := e.(type) {
	case *ContainsExpression:
		return json.Marshal(typedExpression[*ContainsExpression]{
			Type:       "Contains",
			Expression: expr,
		})
	case *IsNotExpression:
		return json.Marshal(typedExpression[*IsNotExpression]{
			Type:       "IsNot",
			Expression: expr,
		})
	case *IsExpression:
		return json.Marshal(typedExpression[*IsExpression]{
			Type:       "Is",
			Expression: expr,
		})
	case *AndExpression:
		return json.Marshal(typedExpression[*AndExpression]{
			Type:       "And",
			Expression: expr,
		})
	case *OrExpression:
		return json.Marshal(typedExpression[*OrExpression]{
			Type:       "Or",
			Expression: expr,
		})
	case *NotExpression:
		return json.Marshal(typedExpression[*NotExpression]{
			Type:       "Not",
			Expression: expr,
		})
	case *GreaterThanExpression:
		return json.Marshal(typedExpression[*GreaterThanExpression]{
			Type:       "GT",
			Expression: expr,
		})
	case *GreaterThanOrEqualExpression:
		return json.Marshal(typedExpression[*GreaterThanOrEqualExpression]{
			Type:       "GTE",
			Expression: expr,
		})
	case *LessThanExpression:
		return json.Marshal(typedExpression[*LessThanExpression]{
			Type:       "LT",
			Expression: expr,
		})
	case *LessThanOrEqualExpression:
		return json.Marshal(typedExpression[*LessThanOrEqualExpression]{
			Type:       "LTE",
			Expression: expr,
		})
	default:
		return nil, fmt.Errorf("unknown expression type %T", e)
	}
}

// unmarshalExpression decodes json data containing a typedExpression and
// returns the underlying Expression.
func unmarshalExpression(data []byte) (Expression, error) {
	var hdr struct{ Type string }
	if err := json.Unmarshal(data, &hdr); err != nil {
		return nil, err
	}
	switch hdr.Type {
	case "Contains":
		var te typedExpression[*ContainsExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "IsNot":
		var te typedExpression[*IsNotExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "Is":
		var te typedExpression[*IsExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "And":
		var te typedExpression[*AndExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "Or":
		var te typedExpression[*OrExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "Not":
		var te typedExpression[*NotExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "GT":
		var te typedExpression[*GreaterThanExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "GTE":
		var te typedExpression[*GreaterThanOrEqualExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "LT":
		var te typedExpression[*LessThanExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	case "LTE":
		var te typedExpression[*LessThanOrEqualExpression]
		if err := json.Unmarshal(data, &te); err != nil {
			return nil, err
		}
		return te.Expression, nil
	default:
		return nil, fmt.Errorf("unrecognized type value %q", hdr.Type)
	}
}

func (q *Query) Evaluate(i interface{}) bool {
	if q.Expression != nil {
		return q.Expression.Evaluate(i)
	}
	return false
}

func (q *Query) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*QueryRaw)(q)); err != nil {
		return err
	}
	if len(q.ExpressionRawJson) == 0 {
		return nil
	}
	expr, err := unmarshalExpression(q.ExpressionRawJson)
	if err != nil {
		return err
	}
	q.Expression = expr
	return nil
}

func (q Query) MarshalJSON() ([]byte, error) {
	if q.Expression != nil {
		data, err := marshalExpression(q.Expression)
		if err != nil {
			return nil, err
		}
		return json.Marshal(&QueryRaw{ExpressionRawJson: data})
	}
	return json.Marshal(&QueryRaw{ExpressionRawJson: q.ExpressionRawJson})
}
