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

// Response は Lambda レスポンスの構造体
type Response struct {
	Output string  `json:"output"`
	Result *string `json:"result"`
	Error  *string `json:"error"`
}

// Handler は Lambda ハンドラー
// event は任意の JSON で、Sugu の event 変数として渡される
func Handler(ctx context.Context, event json.RawMessage) (Response, error) {
	// main.sugu を読み込む
	code, err := os.ReadFile("main.sugu")
	if err != nil {
		errMsg := "failed to read main.sugu: " + err.Error()
		return Response{
			Output: "",
			Result: nil,
			Error:  &errMsg,
		}, nil
	}

	return Execute(string(code), event)
}

// Execute は Sugu コードを実行し、結果を返す
func Execute(code string, eventJSON json.RawMessage) (Response, error) {
	// 出力キャプチャを作成
	capture := &OutputCapture{}

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
		errMsg := "failed to parse event: " + err.Error()
		return Response{
			Output: "",
			Result: nil,
			Error:  &errMsg,
		}, nil
	}
	env.Set("event", eventObj)

	// Lexer
	l := lexer.New(code)

	// Parser
	p := parser.New(l)
	program := p.ParseProgram()

	// パースエラーのチェック
	if len(p.Errors()) != 0 {
		errMsg := p.Errors()[0]
		return Response{
			Output: "",
			Result: nil,
			Error:  &errMsg,
		}, nil
	}

	// Evaluator
	result := evaluator.Eval(program, env)

	// エラーチェック
	if errObj, ok := result.(*object.Error); ok {
		errMsg := errObj.Message
		return Response{
			Output: capture.String(),
			Result: nil,
			Error:  &errMsg,
		}, nil
	}

	// 結果を文字列に変換
	var resultStr *string
	if result != nil {
		s := result.Inspect()
		resultStr = &s
	}

	return Response{
		Output: capture.String(),
		Result: resultStr,
		Error:  nil,
	}, nil
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
