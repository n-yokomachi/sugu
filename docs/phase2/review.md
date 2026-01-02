# Phase 2 レビュー

## PR #12: Add array literal and index expression support

### 概要
配列リテラルとインデックスアクセスの実装。

### 将来の改善点

#### 1. 負のインデックスによる末尾からのアクセス未対応
```go
// evaluator.go:623-624
if idx < 0 || idx > max {
    return NULL
}
```
- Python のように `arr[-1]` で末尾要素にアクセスする機能がない
- 現在は負のインデックスは `null` を返す
- 仕様として意図的であれば問題なし

#### 2. 配列の要素代入が未実装
```sugu
mut arr = [1, 2, 3];
arr[0] = 10;  // これは動作しない
```
- 配列要素への代入式が未実装
- Phase 2 の実装計画に含まれているか確認が必要

#### 3. スライス操作が未実装
```sugu
arr[1:3]  // 部分配列の取得
```
- Python や Go スタイルのスライス構文は未対応
- 将来的に検討

---

## PR #13: Add map literal and key access support

### 概要
マップリテラルとキーアクセスの実装。

### 将来の改善点

#### 1. 空マップと空ブロックの曖昧さ
```sugu
{}  // 空マップ？空ブロック？
```
- 現在の実装では `{}` は空マップとしてパースされる
- ブロック文として解釈されるべきコンテキストでの挙動を確認する必要がある

#### 2. マップの要素代入が未実装
```sugu
mut m = {"a": 1};
m["a"] = 10;  // これは動作しない
m["b"] = 20;  // 新規キーの追加も未対応
```
- マップ要素への代入式が未実装
- 配列と同様に、インデックス代入のサポートが必要

#### 3. 浮動小数点キーのハッシュ衝突
```go
// object/object.go:169-170
func (n *Number) HashKey() HashKey {
    return HashKey{Type: n.Type(), Value: math.Float64bits(n.Value)}
}
```
- `0.0` と `-0.0` が異なるハッシュになる
- `NaN` 同士も異なるハッシュになる可能性
- 浮動小数点をキーとして使用する際の注意が必要

#### 4. マップのキー削除が未実装
```sugu
mut m = {"a": 1, "b": 2};
delete(m, "a");  // キー削除機能がない
```
- 組み込み関数でのキー削除が未対応
- Phase 2 Step 6 の組み込み関数拡張で検討

---

## PR #14: Add builtin functions for array and map operations

### 概要
配列・マップ操作用の組み込み関数を実装。

### 将来の改善点

#### 1. concat() / append() 関数が未実装
```sugu
mut arr1 = [1, 2];
mut arr2 = [3, 4];
concat(arr1, arr2);  // 配列の結合機能がない
```
- 複数配列の結合は現在できない
- `+` 演算子での結合も未対応

#### 2. contains() / includes() 関数が未実装
```sugu
mut arr = [1, 2, 3];
contains(arr, 2);  // 要素の存在確認機能がない
```
- 配列・マップの要素存在確認が組み込みでできない

#### 3. keys()/values() の順序が不定
```go
// builtins.go:145-147
for _, pair := range mapObj.Pairs {
    keys = append(keys, pair.Key)
}
```
- Go の map イテレーションと同様に順序が保証されない
- テストでは len() で要素数のみ確認する工夫が必要

---

## PR #15: Add position information (line/column) to error messages

### 概要
エラーメッセージに行番号・列番号を追加し、デバッグを容易にする。

### 将来の改善点

#### 1. 一部のEvaluatorエラーに位置情報が未適用
```go
// evaluator.go の演算子エラー
return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
return newError("division by zero")
```
- 演算子関連のエラー（型の不一致、ゼロ除算など）には位置情報が含まれていない
- これらのエラーは現在ASTノードを受け取っていないため、リファクタリングが必要
- 優先度は低いが、将来的に対応を検討

#### 2. 組み込み関数のエラーに位置情報がない
```go
// builtins.go
return newError("wrong number of arguments. got=%d, want=1", len(args))
return newError("argument to `len` not supported, got %s", args[0].Type())
```
- 組み込み関数はASTノードにアクセスできないため、位置情報を持てない
- 呼び出し元のCallExpressionから位置情報を渡す仕組みが必要

#### 3. コメント内での改行カウント
```go
// lexer.go:32-38
if l.ch == '\n' {
    l.line++
    l.column = 0
}
```
- 複数行コメント `//-- ... --//` 内の改行も行番号としてカウントされる
- 現在の挙動は正しいが、コメント終了後のトークン位置が正確であることを確認する必要がある

#### 4. タブ文字の列カウント
```go
// lexer.go:37
l.column++
```
- タブ文字も1列としてカウントされる
- エディタによっては4〜8スペースとして表示されるため、列番号がずれる可能性
- 仕様として問題なし（多くの言語で同様の実装）

#### 5. ASTノードへの位置情報の直接保持
```go
// 現在はTokenから位置情報を取得
node.Token.Line, node.Token.Column
```
- 各ASTノードにToken経由でアクセスしているが、直接Line/Columnフィールドを持つ方が明確
- 現在の実装でも動作に問題はない
