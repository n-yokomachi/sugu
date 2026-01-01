package object

import "testing"

func TestEnvironmentGetSet(t *testing.T) {
	env := NewEnvironment()

	// 変数が存在しない場合
	_, ok := env.Get("x")
	if ok {
		t.Error("expected variable 'x' to not exist")
	}

	// 変数を設定
	env.Set("x", &Number{Value: 10})
	obj, ok := env.Get("x")
	if !ok {
		t.Error("expected variable 'x' to exist")
	}
	num, ok := obj.(*Number)
	if !ok {
		t.Errorf("expected Number, got %T", obj)
	}
	if num.Value != 10 {
		t.Errorf("expected value 10, got %f", num.Value)
	}
}

func TestEnvironmentConst(t *testing.T) {
	env := NewEnvironment()

	// const変数を設定
	env.SetConst("PI", &Number{Value: 3.14})

	if !env.IsConst("PI") {
		t.Error("expected 'PI' to be const")
	}

	// mut変数を設定
	env.Set("x", &Number{Value: 10})

	if env.IsConst("x") {
		t.Error("expected 'x' to not be const")
	}
}

func TestEnclosedEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Number{Value: 10})
	outer.SetConst("PI", &Number{Value: 3.14})

	inner := NewEnclosedEnvironment(outer)
	inner.Set("y", &Number{Value: 20})

	// 内部スコープから外部変数にアクセス
	obj, ok := inner.Get("x")
	if !ok {
		t.Error("expected to access outer variable 'x'")
	}
	num := obj.(*Number)
	if num.Value != 10 {
		t.Errorf("expected 10, got %f", num.Value)
	}

	// 内部スコープの変数にアクセス
	obj, ok = inner.Get("y")
	if !ok {
		t.Error("expected to access inner variable 'y'")
	}
	num = obj.(*Number)
	if num.Value != 20 {
		t.Errorf("expected 20, got %f", num.Value)
	}

	// 外部スコープから内部変数にはアクセスできない
	_, ok = outer.Get("y")
	if ok {
		t.Error("expected 'y' to not exist in outer scope")
	}

	// const属性は外部スコープからも参照可能
	if !inner.IsConst("PI") {
		t.Error("expected 'PI' to be const in inner scope")
	}
}

func TestEnvironmentExists(t *testing.T) {
	env := NewEnvironment()
	env.Set("x", &Number{Value: 10})

	if !env.Exists("x") {
		t.Error("expected 'x' to exist")
	}

	if env.Exists("y") {
		t.Error("expected 'y' to not exist")
	}
}

func TestEnvironmentShadowing(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Number{Value: 10})

	inner := NewEnclosedEnvironment(outer)
	inner.Set("x", &Number{Value: 20})

	// 内部スコープでは内部の値が優先
	obj, _ := inner.Get("x")
	num := obj.(*Number)
	if num.Value != 20 {
		t.Errorf("expected 20 in inner scope, got %f", num.Value)
	}

	// 外部スコープでは外部の値のまま
	obj, _ = outer.Get("x")
	num = obj.(*Number)
	if num.Value != 10 {
		t.Errorf("expected 10 in outer scope, got %f", num.Value)
	}
}
