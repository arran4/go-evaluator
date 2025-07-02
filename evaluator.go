package evaluator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Expression interface {
	Evaluate(i interface{}) bool
}

type ContainsExpression struct {
	Field string
	Value interface{}
}

func (e ContainsExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
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

type IsNotExpression struct {
	Field string
	Value interface{}
}

func (e IsNotExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
		return false
	}
	return !reflect.DeepEqual(f.Interface(), e.Value)
}

type IsExpression struct {
	Field string
	Value interface{}
}

func (e IsExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
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

type GreaterThanExpression struct {
	Field string
	Value interface{}
}

func (e GreaterThanExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
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

type GreaterThanOrEqualExpression struct {
	Field string
	Value interface{}
}

func (e GreaterThanOrEqualExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
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

type LessThanExpression struct {
	Field string
	Value interface{}
}

func (e LessThanExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
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

type LessThanOrEqualExpression struct {
	Field string
	Value interface{}
}

func (e LessThanOrEqualExpression) Evaluate(i interface{}) bool {
	v := reflect.ValueOf(i)
	f := v.Elem().FieldByName(e.Field)
	if !f.IsValid() {
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

type QueryRaw struct {
	Expression        Expression      `json:"-"`
	ExpressionRawJson json.RawMessage `json:"Expression"`
}

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
