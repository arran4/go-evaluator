package evaluator

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	return !reflect.DeepEqual(f.Interface(), i)
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
	if f.IsNil() && e.Value == nil {
		return true
	}
	return reflect.DeepEqual(f.Interface(), e.Value)
}

type QueryRaw struct {
	Expression        Expression      `json:"-"`
	ExpressionRawJson json.RawMessage `json:"Expression"`
}

type Query QueryRaw

func (q Query) Evaluate(i interface{}) bool {
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
	default:
		err = fmt.Errorf("unrecognized type value %q", typeStruct.Type)
	}
	if err != nil {
		return err
	}
	return nil
}
