package evaluator

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
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
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				// 文字数を返す（バイト数ではなくrune数）
				return &object.Number{Value: float64(len([]rune(arg.Value)))}
			case *object.Array:
				return &object.Number{Value: float64(len(arg.Elements))}
			case *object.Map:
				return &object.Number{Value: float64(len(arg.Pairs))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return NULL
			}
			newElements := make([]object.Object, length-1)
			copy(newElements, arr.Elements[:length-1])
			return &object.Array{Elements: newElements}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"keys": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.MAP_OBJ {
				return newError("argument to `keys` must be MAP, got %s", args[0].Type())
			}
			mapObj := args[0].(*object.Map)
			keys := make([]object.Object, 0, len(mapObj.Pairs))
			for _, pair := range mapObj.Pairs {
				keys = append(keys, pair.Key)
			}
			return &object.Array{Elements: keys}
		},
	},
	"values": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.MAP_OBJ {
				return newError("argument to `values` must be MAP, got %s", args[0].Type())
			}
			mapObj := args[0].(*object.Map)
			values := make([]object.Object, 0, len(mapObj.Pairs))
			for _, pair := range mapObj.Pairs {
				values = append(values, pair.Value)
			}
			return &object.Array{Elements: values}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Number:
				// 小数点以下切り捨て
				return &object.Number{Value: math.Trunc(arg.Value)}
			case *object.String:
				// 文字列から整数に変換
				val, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newError("cannot convert %q to int", arg.Value)
				}
				return &object.Number{Value: math.Trunc(val)}
			case *object.Boolean:
				if arg.Value {
					return &object.Number{Value: 1}
				}
				return &object.Number{Value: 0}
			default:
				return newError("cannot convert %s to int", args[0].Type())
			}
		},
	},
	"float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Number:
				return arg
			case *object.String:
				val, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newError("cannot convert %q to float", arg.Value)
				}
				return &object.Number{Value: val}
			case *object.Boolean:
				if arg.Value {
					return &object.Number{Value: 1.0}
				}
				return &object.Number{Value: 0.0}
			default:
				return newError("cannot convert %s to float", args[0].Type())
			}
		},
	},
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: args[0].Inspect()}
		},
	},
	"bool": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Number:
				if arg.Value == 0 {
					return FALSE
				}
				return TRUE
			case *object.String:
				if arg.Value == "" {
					return FALSE
				}
				return TRUE
			case *object.Boolean:
				return arg
			case *object.Null:
				return FALSE
			case *object.Array:
				if len(arg.Elements) == 0 {
					return FALSE
				}
				return TRUE
			case *object.Map:
				if len(arg.Pairs) == 0 {
					return FALSE
				}
				return TRUE
			default:
				return TRUE
			}
		},
	},
	"readFile": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `readFile` must be STRING, got %s", args[0].Type())
			}
			path := args[0].(*object.String).Value
			content, err := os.ReadFile(path)
			if err != nil {
				return newError("failed to read file %q: %s", path, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"writeFile": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("first argument to `writeFile` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `writeFile` must be STRING, got %s", args[1].Type())
			}
			path := args[0].(*object.String).Value
			content := args[1].(*object.String).Value
			err := os.WriteFile(path, []byte(content), 0644)
			if err != nil {
				return newError("failed to write file %q: %s", path, err.Error())
			}
			return TRUE
		},
	},
	"appendFile": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("first argument to `appendFile` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `appendFile` must be STRING, got %s", args[1].Type())
			}
			path := args[0].(*object.String).Value
			content := args[1].(*object.String).Value
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return newError("failed to open file %q: %s", path, err.Error())
			}
			defer f.Close()
			_, err = f.WriteString(content)
			if err != nil {
				return newError("failed to append to file %q: %s", path, err.Error())
			}
			return TRUE
		},
	},
	"fileExists": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `fileExists` must be STRING, got %s", args[0].Type())
			}
			path := args[0].(*object.String).Value
			info, err := os.Stat(path)
			if err != nil {
				return FALSE
			}
			// ディレクトリの場合はfalse
			if info.IsDir() {
				return FALSE
			}
			return TRUE
		},
	},
}
