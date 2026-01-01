# Phase 2 レビュー

## PR #12: Add array literal and index expression support

### 概要
配列リテラルとインデックスアクセスの実装。

### 将来の改善点

#### 1. 文字列のインデックスアクセスがバイト単位
```go
// evaluator.go:640
return &object.String{Value: string(stringObject.Value[idx])}
```
- 現在の実装は ASCII 文字のみ正しく動作する
- マルチバイト文字（日本語など）の場合、`[]rune` への変換が必要
- 例: `"あいう"[1]` が期待通りに動作しない可能性

#### 2. 負のインデックスによる末尾からのアクセス未対応
```go
// evaluator.go:623-624
if idx < 0 || idx > max {
    return NULL
}
```
- Python のように `arr[-1]` で末尾要素にアクセスする機能がない
- 現在は負のインデックスは `null` を返す
- 仕様として意図的であれば問題なし

#### 3. 配列の要素代入が未実装
```sugu
mut arr = [1, 2, 3];
arr[0] = 10;  // これは動作しない
```
- 配列要素への代入式が未実装
- Phase 2 の実装計画に含まれているか確認が必要

#### 4. スライス操作が未実装
```sugu
arr[1:3]  // 部分配列の取得
```
- Python や Go スタイルのスライス構文は未対応
- 将来的に検討

#### 5. len() 関数が未実装
```sugu
len([1, 2, 3])  // Error: identifier not found: len
len("hello")   // 同様にエラー
```
- 配列と文字列の長さを取得する `len()` 関数が未実装
- Phase 2 Step 2 の組み込み関数拡張で対応予定

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

#### 3. マップの Inspect() 出力順序が不定
```go
// object/object.go:160-165
for _, pair := range m.Pairs {
    pairs = append(pairs, ...)
}
```
- Go の map イテレーションは順序が保証されない
- テストでの比較やデバッグ時に不便な可能性
- ソートするか、順序を保持する実装を検討

#### 4. 浮動小数点キーのハッシュ衝突
```go
// object/object.go:169-170
func (n *Number) HashKey() HashKey {
    return HashKey{Type: n.Type(), Value: math.Float64bits(n.Value)}
}
```
- `0.0` と `-0.0` が異なるハッシュになる
- `NaN` 同士も異なるハッシュになる可能性
- 浮動小数点をキーとして使用する際の注意が必要

#### 5. マップのキー削除が未実装
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

#### 1. len() の文字列長がバイト単位
```go
// builtins.go:59-60
case *object.String:
    return &object.Number{Value: float64(len(arg.Value))}
```
- `len("あいう")` は 9 を返す（UTF-8 バイト数）
- 文字数を返すには `[]rune` への変換が必要
- Lexer の Unicode 対応と合わせて検討

#### 2. pop() 関数が未実装
```sugu
mut arr = [1, 2, 3];
pop(arr);  // 末尾要素を削除して返す機能がない
```
- push() はあるが pop() がない
- 将来的に追加を検討

#### 3. concat() / append() 関数が未実装
```sugu
mut arr1 = [1, 2];
mut arr2 = [3, 4];
concat(arr1, arr2);  // 配列の結合機能がない
```
- 複数配列の結合は現在できない
- `+` 演算子での結合も未対応

#### 4. contains() / includes() 関数が未実装
```sugu
mut arr = [1, 2, 3];
contains(arr, 2);  // 要素の存在確認機能がない
```
- 配列・マップの要素存在確認が組み込みでできない

#### 5. keys()/values() の順序が不定
```go
// builtins.go:145-147
for _, pair := range mapObj.Pairs {
    keys = append(keys, pair.Key)
}
```
- Go の map イテレーションと同様に順序が保証されない
- テストでは len() で要素数のみ確認する工夫が必要
