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

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "3" {
		t.Errorf("expected result '3', got '%s'", *resp.Result)
	}
}

func TestExecute_OutputCapture(t *testing.T) {
	code := `outln("Hello, Lambda!");`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	expected := "Hello, Lambda!\n"
	if resp.Output != expected {
		t.Errorf("expected output %q, got %q", expected, resp.Output)
	}
}

func TestExecute_MultipleOutputs(t *testing.T) {
	code := `
		outln("Line 1");
		outln("Line 2");
		out("No newline");
	`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	expected := "Line 1\nLine 2\nNo newline"
	if resp.Output != expected {
		t.Errorf("expected output %q, got %q", expected, resp.Output)
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

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "30" {
		t.Errorf("expected result '30', got '%s'", *resp.Result)
	}
}

func TestExecute_Function(t *testing.T) {
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

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "12" {
		t.Errorf("expected result '12', got '%s'", *resp.Result)
	}
}

func TestExecute_ParseError(t *testing.T) {
	code := "1 +"
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response, got nil")
	}

	if resp.Result != nil {
		t.Errorf("expected nil result, got '%s'", *resp.Result)
	}
}

func TestExecute_RuntimeError(t *testing.T) {
	code := "foo"
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response, got nil")
	}

	if resp.Result != nil {
		t.Errorf("expected nil result, got '%s'", *resp.Result)
	}
}

func TestExecute_InFunctionDisabled(t *testing.T) {
	code := `in()`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response, got nil")
	}

	expectedError := "in() is not available in Lambda environment"
	if *resp.Error != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, *resp.Error)
	}
}

func TestExecute_ComplexProgram(t *testing.T) {
	code := `
		func factorial(n) => {
			if (n <= 1) {
				return 1;
			}
			return n * factorial(n - 1);
		}
		outln(factorial(5));
		factorial(5);
	`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	expectedOutput := "120\n"
	if resp.Output != expectedOutput {
		t.Errorf("expected output %q, got %q", expectedOutput, resp.Output)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "120" {
		t.Errorf("expected result '120', got '%s'", *resp.Result)
	}
}

func TestExecute_ArrayAndLoop(t *testing.T) {
	code := `
		mut arr = [1, 2, 3, 4, 5];
		mut sum = 0;
		for (mut i = 0; i < len(arr); i = i + 1) {
			sum = sum + arr[i];
		}
		outln(sum);
		sum;
	`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	expectedOutput := "15\n"
	if resp.Output != expectedOutput {
		t.Errorf("expected output %q, got %q", expectedOutput, resp.Output)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "15" {
		t.Errorf("expected result '15', got '%s'", *resp.Result)
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

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "Taro" {
		t.Errorf("expected result 'Taro', got '%s'", *resp.Result)
	}
}

func TestExecute_EventWithNumber(t *testing.T) {
	code := `event["age"] + 5`
	eventJSON := json.RawMessage(`{"name": "Taro", "age": 25}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "30" {
		t.Errorf("expected result '30', got '%s'", *resp.Result)
	}
}

func TestExecute_EventWithArray(t *testing.T) {
	code := `event["items"][1]`
	eventJSON := json.RawMessage(`{"items": ["apple", "banana", "cherry"]}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "banana" {
		t.Errorf("expected result 'banana', got '%s'", *resp.Result)
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

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "yes" {
		t.Errorf("expected result 'yes', got '%s'", *resp.Result)
	}
}

func TestExecute_EventNull(t *testing.T) {
	code := `event`
	resp, err := Execute(code, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "null" {
		t.Errorf("expected result 'null', got '%s'", *resp.Result)
	}
}

func TestExecute_EventNestedObject(t *testing.T) {
	code := `event["user"]["name"]`
	eventJSON := json.RawMessage(`{"user": {"name": "Hanako", "email": "hanako@example.com"}}`)

	resp, err := Execute(code, eventJSON)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error in response: %s", *resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}

	if *resp.Result != "Hanako" {
		t.Errorf("expected result 'Hanako', got '%s'", *resp.Result)
	}
}
