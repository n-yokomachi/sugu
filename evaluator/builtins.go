package evaluator

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `split` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `split` must be STRING, got %s", args[1].Type())
			}
			str := args[0].(*object.String).Value
			sep := args[1].(*object.String).Value
			parts := strings.Split(str, sep)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `join` must be ARRAY, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `join` must be STRING, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			sep := args[1].(*object.String).Value
			strs := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strs[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(strs, sep)}
		},
	},
	"trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `trim` must be STRING, got %s", args[0].Type())
			}
			str := args[0].(*object.String).Value
			return &object.String{Value: strings.TrimSpace(str)}
		},
	},
	"replace": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `replace` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `replace` must be STRING, got %s", args[1].Type())
			}
			if args[2].Type() != object.STRING_OBJ {
				return newError("third argument to `replace` must be STRING, got %s", args[2].Type())
			}
			str := args[0].(*object.String).Value
			old := args[1].(*object.String).Value
			newStr := args[2].(*object.String).Value
			return &object.String{Value: strings.ReplaceAll(str, old, newStr)}
		},
	},
	"substring": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("first argument to `substring` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.NUMBER_OBJ {
				return newError("second argument to `substring` must be NUMBER, got %s", args[1].Type())
			}
			if args[2].Type() != object.NUMBER_OBJ {
				return newError("third argument to `substring` must be NUMBER, got %s", args[2].Type())
			}
			runes := []rune(args[0].(*object.String).Value)
			start := int(args[1].(*object.Number).Value)
			end := int(args[2].(*object.Number).Value)
			runeLen := len(runes)
			if start < 0 || start > runeLen {
				return newError("substring start index out of range: %d (length: %d)", start, runeLen)
			}
			if end < 0 || end > runeLen {
				return newError("substring end index out of range: %d (length: %d)", end, runeLen)
			}
			if start > end {
				return newError("substring start index %d is greater than end index %d", start, end)
			}
			return &object.String{Value: string(runes[start:end])}
		},
	},
	"indexOf": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("first argument to `indexOf` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `indexOf` must be STRING, got %s", args[1].Type())
			}
			str := args[0].(*object.String).Value
			substr := args[1].(*object.String).Value
			// byte 位置を取得し、rune 位置に変換
			byteIndex := strings.Index(str, substr)
			if byteIndex == -1 {
				return &object.Number{Value: -1}
			}
			runeIndex := len([]rune(str[:byteIndex]))
			return &object.Number{Value: float64(runeIndex)}
		},
	},
	"toUpper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `toUpper` must be STRING, got %s", args[0].Type())
			}
			str := args[0].(*object.String).Value
			return &object.String{Value: strings.ToUpper(str)}
		},
	},
	"toLower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `toLower` must be STRING, got %s", args[0].Type())
			}
			str := args[0].(*object.String).Value
			return &object.String{Value: strings.ToLower(str)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.NUMBER_OBJ {
				return newError("argument to `abs` must be NUMBER, got %s", args[0].Type())
			}
			val := args[0].(*object.Number).Value
			return &object.Number{Value: math.Abs(val)}
		},
	},
	"floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.NUMBER_OBJ {
				return newError("argument to `floor` must be NUMBER, got %s", args[0].Type())
			}
			val := args[0].(*object.Number).Value
			return &object.Number{Value: math.Floor(val)}
		},
	},
	"ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.NUMBER_OBJ {
				return newError("argument to `ceil` must be NUMBER, got %s", args[0].Type())
			}
			val := args[0].(*object.Number).Value
			return &object.Number{Value: math.Ceil(val)}
		},
	},
	"round": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.NUMBER_OBJ {
				return newError("argument to `round` must be NUMBER, got %s", args[0].Type())
			}
			val := args[0].(*object.Number).Value
			return &object.Number{Value: math.Round(val)}
		},
	},
	"sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.NUMBER_OBJ {
				return newError("argument to `sqrt` must be NUMBER, got %s", args[0].Type())
			}
			val := args[0].(*object.Number).Value
			if val < 0 {
				return newError("cannot calculate square root of negative number: %g", val)
			}
			return &object.Number{Value: math.Sqrt(val)}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.NUMBER_OBJ {
				return newError("first argument to `pow` must be NUMBER, got %s", args[0].Type())
			}
			if args[1].Type() != object.NUMBER_OBJ {
				return newError("second argument to `pow` must be NUMBER, got %s", args[1].Type())
			}
			base := args[0].(*object.Number).Value
			exp := args[1].(*object.Number).Value
			return &object.Number{Value: math.Pow(base, exp)}
		},
	},
	"min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=0, want=1+")
			}
			for i, arg := range args {
				if arg.Type() != object.NUMBER_OBJ {
					return newError("argument %d to `min` must be NUMBER, got %s", i+1, arg.Type())
				}
			}
			result := args[0].(*object.Number).Value
			for _, arg := range args[1:] {
				val := arg.(*object.Number).Value
				if val < result {
					result = val
				}
			}
			return &object.Number{Value: result}
		},
	},
	"max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=0, want=1+")
			}
			for i, arg := range args {
				if arg.Type() != object.NUMBER_OBJ {
					return newError("argument %d to `max` must be NUMBER, got %s", i+1, arg.Type())
				}
			}
			result := args[0].(*object.Number).Value
			for _, arg := range args[1:] {
				val := arg.(*object.Number).Value
				if val > result {
					result = val
				}
			}
			return &object.Number{Value: result}
		},
	},
	"random": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.Number{Value: rand.Float64()}
		},
	},
	"delete": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.MAP_OBJ {
				return newError("argument to `delete` must be MAP, got %s", args[0].Type())
			}
			mapObj := args[0].(*object.Map)
			key, ok := args[1].(object.Hashable)
			if !ok {
				return newError("unusable as map key: %s", args[1].Type())
			}
			hashKey := key.HashKey()
			if _, exists := mapObj.Pairs[hashKey]; exists {
				delete(mapObj.Pairs, hashKey)
				return TRUE
			}
			return FALSE
		},
	},
}
