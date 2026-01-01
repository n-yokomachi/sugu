package repl

import (
	"fmt"
	"io"
	"os"
	"sugu/evaluator"
	"sugu/lexer"
	"sugu/object"
	"sugu/parser"
)

// RunFile はファイルを読み込んで実行する
func RunFile(filename string, out io.Writer) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return RunSource(string(content), out)
}

// RunSource はソースコードを実行する
func RunSource(source string, out io.Writer) error {
	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		fmt.Fprintln(out, "Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Fprintf(out, "  %s\n", msg)
		}
		return fmt.Errorf("parse error")
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil {
		if errObj, ok := result.(*object.Error); ok {
			fmt.Fprintf(out, "Error: %s\n", errObj.Message)
			return fmt.Errorf("runtime error: %s", errObj.Message)
		}
	}

	return nil
}
