package evaluator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func derefStructPtr(i interface{}) (reflect.Value, bool) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return reflect.Value{}, false
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	return v, true
}

type Expression interface {
	Evaluate(i interface{}) bool
}

type ContainsExpression struct {
	Field string
	Value interface{}
}

func (e ContainsExpression) Evaluate(i interface{}) bool {
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	v, ok := derefStructPtr(i)
	if !ok {
		return false
	}
	f := v.FieldByName(e.Field)
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
	typeStruct := struct {
		Type string
	}{}
	if err := json.Unmarshal(q.ExpressionRawJson, &typeStruct); err != nil {
		return err
	}
	edata := []byte(q.ExpressionRawJson)
	var err error = nil
	switch typeStruct.Type {
	case "Contains":
		q.Expression = &ContainsExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "IsNot":
		q.Expression = &IsNotExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "Is":
		q.Expression = &IsExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "And":
		q.Expression = &AndExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "Or":
		q.Expression = &OrExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "Not":
		q.Expression = &NotExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "GT":
		q.Expression = &GreaterThanExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "GTE":
		q.Expression = &GreaterThanOrEqualExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "LT":
		q.Expression = &LessThanExpression{}
		err = json.Unmarshal(edata, q.Expression)
	case "LTE":
		q.Expression = &LessThanOrEqualExpression{}
		err = json.Unmarshal(edata, q.Expression)
	default:
		err = fmt.Errorf("unrecognized type value %q", typeStruct.Type)
	}
	if err != nil {
		return err
	}
	return nil
}
