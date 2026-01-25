package evaluator

import "testing"

func TestMapAccess(t *testing.T) {
	// Test map[string]interface{} (Fast Path)
	m1 := map[string]interface{}{
		"Name": "bob",
		"Age":  30,
	}
	if v, err := (IsExpression{Field: "Name", Value: "bob"}.Evaluate(m1)); err != nil || !v {
		t.Errorf("map[string]interface{} access failed: %v %v", v, err)
	}

	// Test map[string]int (Slow Path)
	m2 := map[string]int{
		"Age": 30,
	}
	if v, err := (IsExpression{Field: "Age", Value: 30}.Evaluate(m2)); err != nil || !v {
		t.Errorf("map[string]int access failed: %v %v", v, err)
	}

	// Test map[string]string (Slow Path)
	m3 := map[string]string{
		"Name": "alice",
	}
	if v, err := (IsExpression{Field: "Name", Value: "alice"}.Evaluate(m3)); err != nil || !v {
		t.Errorf("map[string]string access failed: %v %v", v, err)
	}
}

func TestMapNilValue(t *testing.T) {
	m := map[string]interface{}{
		"null": nil,
	}

	// Check if IsExpression handles nil value in map correctly.
	// If getField returns an invalid Value, IsExpression might panic when calling f.Kind().

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic during evaluation: %v", r)
		}
	}()

	expr := IsExpression{Field: "null", Value: nil}

	// This should return true
	if v, err := expr.Evaluate(m); err != nil || !v {
		t.Errorf("IsExpression(nil) failed for nil value in map: %v %v", v, err)
	}
}
