package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"sugu/ast"
	"sugu/object"
)

// シングルトン
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval はASTノードを評価してオブジェクトを返す
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// プログラム
	case *ast.Program:
		return evalProgram(node, env)

	// 文
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.VariableStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if node.Token.Literal == "const" {
			env.SetConst(node.Name.Value, val)
		} else {
			env.Set(node.Name.Value, val)
		}
		return val

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IfStatement:
		return evalIfStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.ForInStatement:
		return evalForInStatement(node, env)

	case *ast.SwitchStatement:
		return evalSwitchStatement(node, env)

	case *ast.BreakStatement:
		return &breakValue{}

	case *ast.ContinueStatement:
		return &continueValue{}

	case *ast.TryStatement:
		return evalTryStatement(node, env)

	case *ast.ThrowStatement:
		return evalThrowStatement(node, env)

	// 式
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.NumberLiteral:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return newErrorWithPos(node.Token.Line, node.Token.Column, "could not parse %q as number", node.Value)
		}
		return &object.Number{Value: val}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return NULL

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		// 論理演算子は短絡評価のため特別扱い
		if node.Operator == "&&" || node.Operator == "||" {
			return evalLogicalExpression(node, env)
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		name := ""
		if node.Name != nil {
			name = node.Name.Value
		}
		fn := &object.Function{Parameters: params, Body: body, Env: env, Name: name}
		if name != "" {
			env.Set(name, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.AssignExpression:
		return evalAssignExpression(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.MapLiteral:
		return evalMapLiteral(node, env)

	case *ast.IndexAssignExpression:
		return evalIndexAssignExpression(node, env)
	}

	return nil
}

// evalProgram はプログラム全体を評価
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		case *throwValue:
			// キャッチされなかったthrowはエラーとして扱う
			return newError("uncaught exception: %s", result.Value.Inspect())
		}
	}

	return result
}

// evalBlockStatement はブロック文を評価
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
			// break/continueの処理
			if _, ok := result.(*breakValue); ok {
				return result
			}
			if _, ok := result.(*continueValue); ok {
				return result
			}
			// throwの処理
			if _, ok := result.(*throwValue); ok {
				return result
			}
		}
	}

	return result
}

// evalIdentifier は識別子を評価
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorWithPos(node.Token.Line, node.Token.Column, "identifier not found: %s", node.Value)
}

// evalAssignExpression は代入式を評価
func evalAssignExpression(node *ast.AssignExpression, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	name := node.Name.Value

	// 変数が存在するか確認
	if _, ok := env.Get(name); !ok {
		return newErrorWithPos(node.Name.Token.Line, node.Name.Token.Column, "identifier not found: %s", name)
	}

	// const変数への再代入をチェック
	if env.IsConst(name) {
		return newErrorWithPos(node.Token.Line, node.Token.Column, "cannot reassign to const variable: %s", name)
	}

	// 変数が定義されているスコープで値を更新
	result, ok := env.Update(name, val)
	if !ok {
		return newErrorWithPos(node.Token.Line, node.Token.Column, "failed to update variable: %s", name)
	}
	return result
}

// evalIndexAssignExpression はインデックス代入式を評価（arr[0] = 10, map["key"] = value）
func evalIndexAssignExpression(node *ast.IndexAssignExpression, env *object.Environment) object.Object {
	// 左辺が識別子の場合、const チェックを行う
	if ident, ok := node.Left.(*ast.Identifier); ok {
		if env.IsConst(ident.Value) {
			return newError("cannot modify const variable: %s", ident.Value)
		}
	}

	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	switch obj := left.(type) {
	case *object.Array:
		return evalArrayIndexAssignment(obj, index, val)
	case *object.Map:
		return evalMapIndexAssignment(obj, index, val)
	default:
		return newError("index assignment not supported: %s", left.Type())
	}
}

// evalArrayIndexAssignment は配列のインデックス代入を評価
func evalArrayIndexAssignment(array *object.Array, index, val object.Object) object.Object {
	if index.Type() != object.NUMBER_OBJ {
		return newError("array index must be a number, got %s", index.Type())
	}

	idx := int64(index.(*object.Number).Value)
	length := int64(len(array.Elements))

	if idx < 0 || idx >= length {
		return newError("array index out of bounds: %d (length: %d)", idx, length)
	}

	array.Elements[idx] = val
	return val
}

// evalMapIndexAssignment はマップのインデックス代入を評価
func evalMapIndexAssignment(mapObj *object.Map, index, val object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	mapObj.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: val}
	return val
}

// evalExpressions は式のリストを評価
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// evalPrefixExpression は前置演算子式を評価
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression は!演算子を評価
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression は-前置演算子を評価
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

// evalInfixExpression は中置演算子式を評価
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalNumberInfixExpression は数値の中置演算子式を評価
func evalNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: leftVal + rightVal}
	case "-":
		return &object.Number{Value: leftVal - rightVal}
	case "*":
		return &object.Number{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Number{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Number{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression は文字列の中置演算子式を評価
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalLogicalExpression は論理演算子を短絡評価する
func evalLogicalExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	if node.Operator == "&&" {
		if !isTruthy(left) {
			return left
		}
		return Eval(node.Right, env)
	}

	// ||
	if isTruthy(left) {
		return left
	}
	return Eval(node.Right, env)
}

// evalIfStatement はif文を評価
func evalIfStatement(ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

// evalWhileStatement はwhile文を評価
func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ws.Body, env)
		if isError(result) {
			return result
		}

		// break処理
		if _, ok := result.(*breakValue); ok {
			return NULL
		}
		// continue処理
		if _, ok := result.(*continueValue); ok {
			continue
		}
		// return処理
		if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}
	}

	return result
}

// evalForStatement はfor文を評価
func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	// 新しいスコープを作成
	forEnv := object.NewEnclosedEnvironment(env)

	// 初期化
	if fs.Init != nil {
		init := Eval(fs.Init, forEnv)
		if isError(init) {
			return init
		}
	}

	var result object.Object = NULL

	for {
		// 条件チェック
		if fs.Condition != nil {
			condition := Eval(fs.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// ボディ実行
		result = Eval(fs.Body, forEnv)
		if isError(result) {
			return result
		}

		// break処理
		if _, ok := result.(*breakValue); ok {
			return NULL
		}
		// continue処理（更新式は実行する）
		if _, ok := result.(*continueValue); ok {
			// 更新式を実行してから続行
			if fs.Update != nil {
				update := Eval(fs.Update, forEnv)
				if isError(update) {
					return update
				}
			}
			continue
		}
		// return処理
		if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}

		// 更新式
		if fs.Update != nil {
			update := Eval(fs.Update, forEnv)
			if isError(update) {
				return update
			}
		}
	}

	return result
}

// evalForInStatement はfor-in文を評価
func evalForInStatement(node *ast.ForInStatement, env *object.Environment) object.Object {
	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	switch obj := iterable.(type) {
	case *object.Array:
		return evalForInArray(node, obj, env)
	case *object.Map:
		return evalForInMap(node, obj, env)
	default:
		return newError("for-in requires ARRAY or MAP, got %s", iterable.Type())
	}
}

// evalForInArray は配列に対する for-in を評価
func evalForInArray(node *ast.ForInStatement, arr *object.Array, env *object.Environment) object.Object {
	var result object.Object = NULL

	for i, elem := range arr.Elements {
		iterEnv := object.NewEnclosedEnvironment(env)

		if node.Value != nil {
			// for (index, item in arr)
			iterEnv.SetConst(node.Key.Value, &object.Number{Value: float64(i)})
			iterEnv.SetConst(node.Value.Value, elem)
		} else {
			// for (item in arr)
			iterEnv.SetConst(node.Key.Value, elem)
		}

		result = Eval(node.Body, iterEnv)
		if isError(result) {
			return result
		}
		if _, ok := result.(*breakValue); ok {
			return NULL
		}
		if _, ok := result.(*continueValue); ok {
			continue
		}
		if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}
		if _, ok := result.(*throwValue); ok {
			return result
		}
	}

	return result
}

// evalForInMap はマップに対する for-in を評価
func evalForInMap(node *ast.ForInStatement, mapObj *object.Map, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, pair := range mapObj.Pairs {
		iterEnv := object.NewEnclosedEnvironment(env)

		if node.Value != nil {
			// for (key, value in map)
			iterEnv.SetConst(node.Key.Value, pair.Key)
			iterEnv.SetConst(node.Value.Value, pair.Value)
		} else {
			// for (key in map)
			iterEnv.SetConst(node.Key.Value, pair.Key)
		}

		result = Eval(node.Body, iterEnv)
		if isError(result) {
			return result
		}
		if _, ok := result.(*breakValue); ok {
			return NULL
		}
		if _, ok := result.(*continueValue); ok {
			continue
		}
		if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}
		if _, ok := result.(*throwValue); ok {
			return result
		}
	}

	return result
}

// evalSwitchStatement はswitch文を評価
func evalSwitchStatement(ss *ast.SwitchStatement, env *object.Environment) object.Object {
	value := Eval(ss.Value, env)
	if isError(value) {
		return value
	}

	for _, caseClause := range ss.Cases {
		caseValue := Eval(caseClause.Value, env)
		if isError(caseValue) {
			return caseValue
		}

		if isEqual(value, caseValue) {
			result := Eval(caseClause.Body, env)
			// breakを無視（switch内では自動break）
			if _, ok := result.(*breakValue); ok {
				return NULL
			}
			return result
		}
	}

	// default節
	if ss.Default != nil {
		result := Eval(ss.Default, env)
		if _, ok := result.(*breakValue); ok {
			return NULL
		}
		return result
	}

	return NULL
}

// applyFunction は関数を適用
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv は関数用の新しい環境を作成
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
		} else {
			env.Set(param.Value, NULL)
		}
	}

	return env
}

// unwrapReturnValue はReturnValueをアンラップ
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// ヘルパー関数

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isEqual(a, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *object.Number:
		return a.Value == b.(*object.Number).Value
	case *object.String:
		return a.Value == b.(*object.String).Value
	case *object.Boolean:
		return a.Value == b.(*object.Boolean).Value
	case *object.Null:
		return true
	default:
		return a == b
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func newErrorWithPos(line, column int, format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	return &object.Error{Message: fmt.Sprintf("line %d, column %d: %s", line, column, msg)}
}

// evalIndexExpression はインデックスアクセスを評価
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.MAP_OBJ:
		return evalMapIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalMapLiteral はマップリテラルを評価
func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

// evalMapIndexExpression はマップのインデックスアクセスを評価
func evalMapIndexExpression(mapObj, index object.Object) object.Object {
	mapObject := mapObj.(*object.Map)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := mapObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

// evalArrayIndexExpression は配列のインデックスアクセスを評価
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := int64(index.(*object.Number).Value)
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

// evalStringIndexExpression は文字列のインデックスアクセスを評価
func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	// runeに変換してマルチバイト文字に対応
	runes := []rune(stringObject.Value)
	idx := int64(index.(*object.Number).Value)
	max := int64(len(runes) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string(runes[idx])}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// break/continue用の内部型
type breakValue struct{}

func (b *breakValue) Type() object.ObjectType { return "BREAK" }
func (b *breakValue) Inspect() string         { return "break" }

type continueValue struct{}

func (c *continueValue) Type() object.ObjectType { return "CONTINUE" }
func (c *continueValue) Inspect() string         { return "continue" }

// throw用の内部型
type throwValue struct {
	Value object.Object
}

func (t *throwValue) Type() object.ObjectType { return "THROW" }
func (t *throwValue) Inspect() string         { return "throw " + t.Value.Inspect() }

// evalTryStatement はtry/catch文を評価
func evalTryStatement(ts *ast.TryStatement, env *object.Environment) object.Object {
	// tryブロックを評価
	result := Eval(ts.TryBlock, env)

	// throwされた場合、catchブロックを実行
	if thrown, ok := result.(*throwValue); ok {
		// 新しいスコープでcatchブロックを実行
		catchEnv := object.NewEnclosedEnvironment(env)
		catchEnv.Set(ts.CatchParam.Value, thrown.Value)
		return Eval(ts.CatchBlock, catchEnv)
	}

	// Errorオブジェクト（組み込み関数などからのエラー）もキャッチ
	if errObj, ok := result.(*object.Error); ok {
		catchEnv := object.NewEnclosedEnvironment(env)
		catchEnv.Set(ts.CatchParam.Value, &object.String{Value: errObj.Message})
		return Eval(ts.CatchBlock, catchEnv)
	}

	return result
}

// evalThrowStatement はthrow文を評価
func evalThrowStatement(ts *ast.ThrowStatement, env *object.Environment) object.Object {
	val := Eval(ts.Value, env)
	if isError(val) {
		return val
	}
	return &throwValue{Value: val}
}
