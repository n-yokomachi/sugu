package repl

import (
	"bufio"
	"fmt"
	"io"
	"sugu/evaluator"
	"sugu/lexer"
	"sugu/object"
	"sugu/parser"
)

const PROMPT = ">> "

// Start はREPLを開始する
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "Sugu Language REPL")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit")
	fmt.Fprintln(out)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// 終了コマンド
		if line == "exit" || line == "quit" {
			fmt.Fprintln(out, "Bye!")
			return
		}

		// 空行はスキップ
		if line == "" {
			continue
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if errObj, ok := evaluated.(*object.Error); ok {
				fmt.Fprintf(out, "Error: %s\n", errObj.Message)
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

// printParserErrors はパーサーエラーを表示する
func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, "Parser errors:")
	for _, msg := range errors {
		fmt.Fprintf(out, "  %s\n", msg)
	}
}
