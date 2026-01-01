package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPLBasicExpression(t *testing.T) {
	input := "1 + 2\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "3") {
		t.Errorf("expected output to contain '3', got=%q", output)
	}
	if !strings.Contains(output, "Bye!") {
		t.Errorf("expected output to contain 'Bye!', got=%q", output)
	}
}

func TestREPLVariable(t *testing.T) {
	input := "mut x = 10;\nx * 2\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "20") {
		t.Errorf("expected output to contain '20', got=%q", output)
	}
}

func TestREPLFunction(t *testing.T) {
	input := "func add(a, b) => { a + b; };\nadd(3, 4)\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "7") {
		t.Errorf("expected output to contain '7', got=%q", output)
	}
}

func TestREPLParserError(t *testing.T) {
	input := "1 +\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "Parser errors") {
		t.Errorf("expected output to contain 'Parser errors', got=%q", output)
	}
}

func TestREPLRuntimeError(t *testing.T) {
	input := "x\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("expected output to contain 'Error:', got=%q", output)
	}
	if !strings.Contains(output, "identifier not found") {
		t.Errorf("expected output to contain 'identifier not found', got=%q", output)
	}
}

func TestREPLQuit(t *testing.T) {
	input := "quit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "Bye!") {
		t.Errorf("expected output to contain 'Bye!', got=%q", output)
	}
}

func TestREPLEmptyLine(t *testing.T) {
	input := "\n\n1\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "1") {
		t.Errorf("expected output to contain '1', got=%q", output)
	}
}

func TestRunSource(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
	}{
		{
			name:        "simple expression",
			source:      "1 + 2",
			expectError: false,
		},
		{
			name:        "variable declaration",
			source:      "mut x = 10; x * 2;",
			expectError: false,
		},
		{
			name:        "function definition and call",
			source:      "func add(a, b) => { a + b; }; add(3, 4);",
			expectError: false,
		},
		{
			name:        "control flow",
			source:      "mut sum = 0; for (mut i = 1; i <= 5; i = i + 1) { sum = sum + i; } sum;",
			expectError: false,
		},
		{
			name:        "parser error",
			source:      "1 +",
			expectError: true,
		},
		{
			name:        "runtime error - undefined variable",
			source:      "x",
			expectError: true,
		},
		{
			name:        "runtime error - division by zero",
			source:      "10 / 0",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := RunSource(tt.source, out)

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
