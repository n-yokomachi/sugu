package main

import (
	"strings"
	"sugu/object"
)

// OutputCapture は Lambda 環境での出力をキャプチャする
type OutputCapture struct {
	builder strings.Builder
}

// Write は文字列を出力バッファに書き込む
func (o *OutputCapture) Write(s string) {
	o.builder.WriteString(s)
}

// String はキャプチャした出力を返す
func (o *OutputCapture) String() string {
	return o.builder.String()
}

// Reset はバッファをクリアする
func (o *OutputCapture) Reset() {
	o.builder.Reset()
}

// NewLambdaBuiltins は Lambda 用の組み込み関数を作成する
func NewLambdaBuiltins(capture *OutputCapture) map[string]*object.Builtin {
	return map[string]*object.Builtin{
		"out": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					capture.Write(arg.Inspect())
				}
				return &object.Null{}
			},
		},
		"outln": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					capture.Write(arg.Inspect())
					capture.Write("\n")
				}
				return &object.Null{}
			},
		},
		"in": {
			Fn: func(args ...object.Object) object.Object {
				return &object.Error{Message: "in() is not available in Lambda environment"}
			},
		},
	}
}
