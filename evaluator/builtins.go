package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"sugu/object"
)

var builtins = map[string]*object.Builtin{
	"out": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	"outln": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"in": {
		Fn: func(args ...object.Object) object.Object {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return newError("failed to read input: %s", err.Error())
			}
			// 末尾の改行を削除
			if len(input) > 0 && input[len(input)-1] == '\n' {
				input = input[:len(input)-1]
			}
			// Windowsの場合、CRも削除
			if len(input) > 0 && input[len(input)-1] == '\r' {
				input = input[:len(input)-1]
			}
			return &object.String{Value: input}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
}
