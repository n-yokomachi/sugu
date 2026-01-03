package main

import (
	"encoding/json"
	"testing"
)

func TestExecute_SimpleExpression(t *testing.T) {
	code := "1 + 2"
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != float64(3) {
		t.Errorf("expected 3, got %v", resp)
	}
}

func TestExecute_ReturnString(t *testing.T) {
	code := `"hello"`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "hello" {
		t.Errorf("expected 'hello', got %v", resp)
	}
}

func TestExecute_ReturnMap(t *testing.T) {
	code := `{"status": "ok", "count": 42}`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", resp)
	}

	if resultMap["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", resultMap["status"])
	}

	if resultMap["count"] != float64(42) {
		t.Errorf("expected count 42, got %v", resultMap["count"])
	}
}

func TestExecute_ReturnArray(t *testing.T) {
	code := `[1, 2, 3]`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultArr, ok := resp.([]interface{})
	if !ok {
		t.Fatalf("expected array, got %T", resp)
	}

	if len(resultArr) != 3 {
		t.Errorf("expected 3 elements, got %d", len(resultArr))
	}

	if resultArr[0] != float64(1) {
		t.Errorf("expected first element 1, got %v", resultArr[0])
	}
}

func TestExecute_VariableAndExpression(t *testing.T) {
	code := `
		mut x = 10;
		mut y = 20;
		x + y;
	`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != float64(30) {
		t.Errorf("expected 30, got %v", resp)
	}
}

func TestExecute_ParseError(t *testing.T) {
	code := "1 +"
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := resp.(map[string]string)
	if !ok {
		t.Fatalf("expected error map, got %T", resp)
	}

	if resultMap["error"] == "" {
		t.Fatal("expected error message")
	}
}

func TestExecute_RuntimeError(t *testing.T) {
	code := "foo"
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := resp.(map[string]string)
	if !ok {
		t.Fatalf("expected error map, got %T", resp)
	}

	if resultMap["error"] == "" {
		t.Fatal("expected error message")
	}
}

func TestExecute_InFunctionDisabled(t *testing.T) {
	code := `in()`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := resp.(map[string]string)
	if !ok {
		t.Fatalf("expected error map, got %T", resp)
	}

	expectedError := "in() is not available in Lambda environment"
	if resultMap["error"] != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, resultMap["error"])
	}
}

func TestExecute_FunctionCall(t *testing.T) {
	code := `
		func add(a, b) => {
			return a + b;
		}
		add(5, 7);
	`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != float64(12) {
		t.Errorf("expected 12, got %v", resp)
	}
}

func TestOutputCapture_Reset(t *testing.T) {
	capture := &OutputCapture{}
	capture.Write("test")

	if capture.String() != "test" {
		t.Errorf("expected 'test', got '%s'", capture.String())
	}

	capture.Reset()

	if capture.String() != "" {
		t.Errorf("expected empty string after reset, got '%s'", capture.String())
	}
}

// Event 変数のテスト

func TestExecute_EventVariable(t *testing.T) {
	code := `event["name"]`
	eventJSON := json.RawMessage(`{"name": "Taro", "age": 25}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "Taro" {
		t.Errorf("expected 'Taro', got %v", resp)
	}
}

func TestExecute_EventWithNumber(t *testing.T) {
	code := `event["age"] + 5`
	eventJSON := json.RawMessage(`{"name": "Taro", "age": 25}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != float64(30) {
		t.Errorf("expected 30, got %v", resp)
	}
}

func TestExecute_EventWithArray(t *testing.T) {
	code := `event["items"][1]`
	eventJSON := json.RawMessage(`{"items": ["apple", "banana", "cherry"]}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "banana" {
		t.Errorf("expected 'banana', got %v", resp)
	}
}

func TestExecute_EventWithBoolean(t *testing.T) {
	code := `
		if (event["active"]) {
			"yes";
		} else {
			"no";
		}
	`
	eventJSON := json.RawMessage(`{"active": true}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "yes" {
		t.Errorf("expected 'yes', got %v", resp)
	}
}

func TestExecute_EventNull(t *testing.T) {
	code := `event`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != nil {
		t.Errorf("expected nil, got %v", resp)
	}
}

func TestExecute_EventNestedObject(t *testing.T) {
	code := `event["user"]["name"]`
	eventJSON := json.RawMessage(`{"user": {"name": "Hanako", "email": "hanako@example.com"}}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "Hanako" {
		t.Errorf("expected 'Hanako', got %v", resp)
	}
}

func TestExecute_ReturnEventAsIs(t *testing.T) {
	code := `event`
	eventJSON := json.RawMessage(`{"name": "Test", "count": 10}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", resp)
	}

	if resultMap["name"] != "Test" {
		t.Errorf("expected name 'Test', got %v", resultMap["name"])
	}

	if resultMap["count"] != float64(10) {
		t.Errorf("expected count 10, got %v", resultMap["count"])
	}
}
