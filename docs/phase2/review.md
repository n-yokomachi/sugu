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
