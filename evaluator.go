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

// number represents any built-in numeric type.
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// numeric converts v into the requested numeric type if possible.
func numeric[T number](v interface{}) (T, bool) {
	var zero T
	switch n := v.(type) {
	case int:
		return T(n), true
	case int8:
		return T(n), true
	case int16:
		return T(n), true
	case int32:
		return T(n), true
	case int64:
		return T(n), true
	case uint:
		return T(n), true
	case uint8:
		return T(n), true
	case uint16:
		return T(n), true
	case uint32:
		return T(n), true
	case uint64:
		return T(n), true
	case uintptr:
		return T(n), true
	case float32:
		return T(n), true
	case float64:
		return T(n), true
	case json.Number:
		f, err := n.Float64()
		if err == nil {
			return T(f), true
		}
		return zero, false
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return T(f), true
		}
		return zero, false
	default:
		// Attempt reflection for other numeric types that might not match exact types in switch
		// e.g. int vs int64 mismatches if T is int64 but v is int
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return T(val.Int()), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return T(val.Uint()), true
		case reflect.Float32, reflect.Float64:
			return T(val.Float()), true
		}
		return zero, false
	}
}

// Comparator allows for custom comparison logic.
type Comparator interface {
	Compare(other interface{}) (int, error)
}

// Compare returns an integer comparing two values.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a, b interface{}) (int, error) {
	if c, ok := a.(Comparator); ok {
		return c.Compare(b)
	}
	if n1, ok := numeric[float64](a); ok {
		if n2, ok := numeric[float64](b); ok {
			if n1 < n2 {
				return -1, nil
			}
			if n1 > n2 {
				return 1, nil
			}
			return 0, nil
		}
	}
	s1 := stringValue(a)
	s2 := stringValue(b)
	return strings.Compare(s1, s2), nil
}

func stringValue(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	default:
		return fmt.Sprint(v)
	}
}

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

// Getter interface allows for dynamic field retrieval.
type Getter interface {
	Get(name string) (interface{}, error)
}

// getField retrieves a field value from either a struct, map, or Getter.
// For structs it uses FieldByName, for maps it looks up the key by name,
// and for Getter it calls Get.
func getField(v reflect.Value, name string) (reflect.Value, bool) {
	if v.Kind() == reflect.Invalid {
		return reflect.Value{}, false
	}
	if v.CanInterface() {
		if g, ok := v.Interface().(Getter); ok {
			val, err := g.Get(name)
			if err == nil {
				if val == nil {
					// Handle nil interface return
					return reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem()), true
				}
				return reflect.ValueOf(val), true
			}
			return reflect.Value{}, false
		}
	}

	switch v.Kind() {
	case reflect.Struct:
		f := v.FieldByName(name)
		if f.IsValid() {
			return f, true
		}
		return reflect.Value{}, false
	case reflect.Map:
		// Fast path for map[string]interface{}
		if v.CanInterface() {
			if m, ok := v.Interface().(map[string]interface{}); ok {
				if val, found := m[name]; found {
					if val == nil {
						return reflect.Zero(v.Type().Elem()), true
					}
					return reflect.ValueOf(val), true
				}
				return reflect.Value{}, false
			}
		}

		key := reflect.ValueOf(name)
		if key.Type().AssignableTo(v.Type().Key()) {
			f := v.MapIndex(key)
			if f.IsValid() {
				// If the map value is an interface, we need to extract the underlying value
				// to get the correct Kind() for comparison.
				if f.Kind() == reflect.Interface {
					return f.Elem(), true
				}
				return f, true
			}
		}
		return reflect.Value{}, false
	default:
		return reflect.Value{}, false
	}
}

func greater[T number](f T, v interface{}) bool {
	n, ok := numeric[T](v)
	if !ok {
		return false
	}
	return f > n
}

func greaterOrEqual[T number](f T, v interface{}) bool {
	n, ok := numeric[T](v)
	if !ok {
		return false
	}
	return f >= n
}

func less[T number](f T, v interface{}) bool {
	n, ok := numeric[T](v)
	if !ok {
		return false
	}
	return f < n
}

func lessOrEqual[T number](f T, v interface{}) bool {
	n, ok := numeric[T](v)
	if !ok {
		return false
	}
	return f <= n
}

type Term interface {
	Evaluate(i interface{}) (interface{}, error)
}

// Function defines the interface for a function that can be called by FunctionExpression.
type Function interface {
	Call(args ...interface{}) (interface{}, error)
}

// FunctionExpression represents a function call.
type FunctionExpression struct {
	Func Function
	Args []Term
}

func (f FunctionExpression) Evaluate(i interface{}) (interface{}, error) {
	args := make([]interface{}, len(f.Args))
	for idx, arg := range f.Args {
		val, err := arg.Evaluate(i)
		if err != nil {
			return nil, err
		}
		args[idx] = val
	}
	return f.Func.Call(args...)
}

// Field represents a field lookup term.
type Field struct {
	Name string
}

func (f Field) Evaluate(i interface{}) (interface{}, error) {
	v, ok := derefValue(i)
	if !ok {
		return nil, fmt.Errorf("cannot dereference value")
	}
	val, ok := getField(v, f.Name)
	if !ok {
		return nil, fmt.Errorf("field %s not found", f.Name)
	}
	if val.IsValid() && val.CanInterface() {
		return val.Interface(), nil
	}
	return nil, nil
}

// Constant represents a constant value term.
type Constant struct {
	Value interface{}
}

func (c Constant) Evaluate(i interface{}) (interface{}, error) {
	return c.Value, nil
}

// Self represents the input value itself.
type Self struct{}

func (s Self) Evaluate(i interface{}) (interface{}, error) {
	return i, nil
}

// BoolType converts the term result to a boolean.
type BoolType struct {
	Term Term
}

func (b BoolType) Evaluate(i interface{}) (interface{}, error) {
	val, err := b.Term.Evaluate(i)
	if err != nil {
		return false, err
	}
	v, err := IsTruthy(val)
	return v, err
}

// IsTruthy checks if a value is considered "true" in the expression language.
// It tries to accept widely accepted truthy values, including parsing strings.
func IsTruthy(v interface{}) (bool, error) {
	if v == nil {
		return false, nil
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0, nil
	case reflect.String:
		b, err := strconv.ParseBool(val.String())
		if err != nil {
			return false, err
		}
		return b, nil
	case reflect.Slice, reflect.Map, reflect.Chan:
		return !val.IsNil() && val.Len() > 0, nil
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return false, nil
		}
		return IsTruthy(val.Elem().Interface())
	}
	return true, nil
}

// If evaluates Condition. If true, evaluates Then, else evaluates Else.
type If struct {
	Condition Term
	Then      Term
	Else      Term
}

func (e If) Evaluate(i interface{}) (interface{}, error) {
	condVal, err := e.Condition.Evaluate(i)
	if err != nil {
		return nil, err
	}
	b, err := IsTruthy(condVal)
	if err != nil {
		return nil, err
	}
	if b {
		return e.Then.Evaluate(i)
	}
	if e.Else != nil {
		return e.Else.Evaluate(i)
	}
	return nil, nil // Or appropriate zero value
}

// ComparisonExpression evaluates a comparison between two Terms.
type ComparisonExpression struct {
	LHS       Term
	RHS       Term
	Operation string // eq, neq, gt, gte, lt, lte, contains, icontains
}

func (e ComparisonExpression) Evaluate(i interface{}) bool {
	lhs, err := e.LHS.Evaluate(i)
	if err != nil {
		return false
	}
	rhs, err := e.RHS.Evaluate(i)
	if err != nil {
		return false
	}

	switch e.Operation {
	case "eq":
		cmp, err := Compare(lhs, rhs)
		return err == nil && cmp == 0
	case "neq":
		cmp, err := Compare(lhs, rhs)
		return err == nil && cmp != 0
	case "gt":
		cmp, err := Compare(lhs, rhs)
		return err == nil && cmp > 0
	case "gte":
		cmp, err := Compare(lhs, rhs)
		return err == nil && cmp >= 0
	case "lt":
		cmp, err := Compare(lhs, rhs)
		return err == nil && cmp < 0
	case "lte":
		cmp, err := Compare(lhs, rhs)
		return err == nil && cmp <= 0
	case "contains":
		s1 := stringValue(lhs)
		s2 := stringValue(rhs)
		return strings.Contains(s1, s2)
	case "icontains":
		s1 := stringValue(lhs)
		s2 := stringValue(rhs)
		return strings.Contains(strings.ToLower(s1), strings.ToLower(s2))
	}
	return false
}

// Expression represents a single boolean expression that can be evaluated
// against a struct value.
type Expression interface {
	// Evaluate returns true if the expression matches the supplied value.
	Evaluate(i interface{}) bool
}

// ContainsExpression checks whether a slice field contains the given Value,
// or if a string field contains the given substring.
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
	if f.Kind() == reflect.String {
		sval := stringValue(e.Value)
		return strings.Contains(f.String(), sval)
	}
	if f.Kind() != reflect.Slice {
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

// IContainsExpression checks whether a string field contains the given substring (case-insensitive).
type IContainsExpression struct {
	Field string
	Value interface{}
}

func (e IContainsExpression) Evaluate(i interface{}) bool {
	v, ok := derefValue(i)
	if !ok {
		return false
	}
	f, ok := getField(v, e.Field)
	if !ok {
		return false
	}
	if f.Kind() == reflect.String {
		sval := stringValue(e.Value)
		return strings.Contains(strings.ToLower(f.String()), strings.ToLower(sval))
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
	if reflect.DeepEqual(f.Interface(), e.Value) {
		return true
	}
	return stringValue(f.Interface()) == stringValue(e.Value)
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
		return greater[int64](f.Int(), e.Value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return greater[uint64](f.Uint(), e.Value)
	case reflect.Float32, reflect.Float64:
		return greater[float64](f.Float(), e.Value)
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
		return greaterOrEqual[int64](f.Int(), e.Value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return greaterOrEqual[uint64](f.Uint(), e.Value)
	case reflect.Float32, reflect.Float64:
		return greaterOrEqual[float64](f.Float(), e.Value)
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
		return less[int64](f.Int(), e.Value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return less[uint64](f.Uint(), e.Value)
	case reflect.Float32, reflect.Float64:
		return less[float64](f.Float(), e.Value)
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
		return lessOrEqual[int64](f.Int(), e.Value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return lessOrEqual[uint64](f.Uint(), e.Value)
	case reflect.Float32, reflect.Float64:
		return lessOrEqual[float64](f.Float(), e.Value)
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
	case *IContainsExpression:
		return json.Marshal(typedExpression[*IContainsExpression]{
			Type:       "IContains",
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
	case "IContains":
		var te typedExpression[*IContainsExpression]
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
