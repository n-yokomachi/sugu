package object

// Environment は変数のスコープを管理する
type Environment struct {
	store  map[string]Object
	consts map[string]bool // const宣言された変数を追跡
	outer  *Environment
}

// NewEnvironment は新しい環境を作成する
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, consts: c, outer: nil}
}

// NewEnclosedEnvironment は外部環境を持つ新しい環境を作成する
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get は変数の値を取得する
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set は変数に値を設定する
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// SetConst はconst変数として値を設定する
func (e *Environment) SetConst(name string, val Object) Object {
	e.store[name] = val
	e.consts[name] = true
	return val
}

// IsConst は変数がconstかどうかを返す
func (e *Environment) IsConst(name string) bool {
	if isConst, ok := e.consts[name]; ok {
		return isConst
	}
	if e.outer != nil {
		return e.outer.IsConst(name)
	}
	return false
}

// Exists は変数が現在のスコープに存在するかを返す
func (e *Environment) Exists(name string) bool {
	_, ok := e.store[name]
	return ok
}

// Update は変数が定義されているスコープで値を更新する
func (e *Environment) Update(name string, val Object) (Object, bool) {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val, true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return nil, false
}
