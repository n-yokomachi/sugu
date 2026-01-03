package main

import (
	"context"
	"encoding/json"
	"os"
	"sugu/evaluator"
	"sugu/lexer"
	"sugu/object"
	"sugu/parser"
)

// Handler は Lambda ハンドラー
// event は任意の JSON で、Sugu の event 変数として渡される
// main.sugu の最後に評価された式の値がそのまま Lambda のレスポンスになる
func Handler(ctx context.Context, event json.RawMessage) (interface{}, error) {
	// main.sugu を読み込む
	code, err := os.ReadFile("main.sugu")
	if err != nil {
		return map[string]string{"error": "failed to read main.sugu: " + err.Error()}, nil
	}

	return Execute(string(code), event)
}

// Execute は Sugu コードを実行し、結果を返す
func Execute(code string, eventJSON json.RawMessage) (interface{}, error) {
	// 出力キャプチャを作成
	capture := &OutputCapture{}
	_ = capture // outln の出力は Lambda では使用しないが、エラー防止のため保持

	// Lambda 用組み込み関数を取得
	lambdaBuiltins := NewLambdaBuiltins(capture)

	// 環境を作成し、Lambda 用組み込み関数を登録
	env := object.NewEnvironment()
	for name, builtin := range lambdaBuiltins {
		env.Set(name, builtin)
	}

	// event 変数を設定
	eventObj, err := jsonToSuguObject(eventJSON)
	if err != nil {
		return map[string]string{"error": "failed to parse event: " + err.Error()}, nil
	}
	env.Set("event", eventObj)

	// Lexer
	l := lexer.New(code)

	// Parser
	p := parser.New(l)
	program := p.ParseProgram()

	// パースエラーのチェック
	if len(p.Errors()) != 0 {
		return map[string]string{"error": p.Errors()[0]}, nil
	}

	// Evaluator
	result := evaluator.Eval(program, env)

	// エラーチェック
	if errObj, ok := result.(*object.Error); ok {
		return map[string]string{"error": errObj.Message}, nil
	}

	// 結果を返す（return 文の値、または最後に評価された式の値）
	return suguObjectToGoValue(result), nil
}

// jsonToSuguObject は JSON を Sugu の Object に変換する
func jsonToSuguObject(data json.RawMessage) (object.Object, error) {
	if len(data) == 0 {
		return &object.Null{}, nil
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return goValueToSuguObject(v), nil
}

// goValueToSuguObject は Go の値を Sugu の Object に変換する
func goValueToSuguObject(v interface{}) object.Object {
	switch val := v.(type) {
	case nil:
		return &object.Null{}
	case bool:
		return &object.Boolean{Value: val}
	case float64:
		return &object.Number{Value: val}
	case string:
		return &object.String{Value: val}
	case []interface{}:
		elements := make([]object.Object, len(val))
		for i, elem := range val {
			elements[i] = goValueToSuguObject(elem)
		}
		return &object.Array{Elements: elements}
	case map[string]interface{}:
		pairs := make(map[object.HashKey]object.HashPair)
		for k, v := range val {
			key := &object.String{Value: k}
			pairs[key.HashKey()] = object.HashPair{
				Key:   key,
				Value: goValueToSuguObject(v),
			}
		}
		return &object.Map{Pairs: pairs}
	default:
		return &object.Null{}
	}
}

// suguObjectToGoValue は Sugu の Object を Go の値に変換する
func suguObjectToGoValue(obj object.Object) interface{} {
	switch v := obj.(type) {
	case *object.Null:
		return nil
	case *object.Boolean:
		return v.Value
	case *object.Number:
		return v.Value
	case *object.String:
		return v.Value
	case *object.Array:
		result := make([]interface{}, len(v.Elements))
		for i, elem := range v.Elements {
			result[i] = suguObjectToGoValue(elem)
		}
		return result
	case *object.Map:
		result := make(map[string]interface{})
		for _, pair := range v.Pairs {
			if key, ok := pair.Key.(*object.String); ok {
				result[key.Value] = suguObjectToGoValue(pair.Value)
			}
		}
		return result
	default:
		return nil
	}
}
