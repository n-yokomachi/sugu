# Requirements Document

## Introduction

cc-sdd 導入前に行われた Sugu 言語の実装状況を棚卸しし、次のフェーズ（Phase 4）の実装要件を定義する。Phase 1〜3 と Lambda 対応は完了済み。Phase 4 実装計画（`docs/phase4/implementation.md`）は既に存在するが、cc-sdd のワークフローに乗せて仕様駆動で進めるための要件として再定義する。

### 現在の実装状況

| Phase | 内容 | 状態 |
|-------|------|------|
| Phase 1 | 基本構文（Lexer, Parser, Evaluator, REPL） | 完了 |
| Phase 2 | 配列・マップ・組み込み関数・位置情報 | 完了 |
| Phase 3 | 代入式・try/catch/throw・型変換・ファイルI/O | 完了 |
| Lambda | AWS Lambda 対応 | 完了 |
| Phase 4 | 文字列操作・数学関数・for-in 等 | **未着手** |

全テスト通過済み（ast, evaluator, lambda, lexer, object, parser, repl）。

---

## Requirements

### Requirement 1: 剰余演算子の浮動小数点対応

**Objective:** 開発者として、剰余演算子 `%` が浮動小数点数でも正しく動作してほしい。これにより数値計算の正確性が向上する。

#### Acceptance Criteria

1. When ユーザーが浮動小数点数同士の `%` 演算を実行した場合, the Sugu Evaluator shall `math.Mod` を使用して正しい浮動小数点剰余を返す（例: `5.5 % 2.0` → `1.5`）
2. When ユーザーが整数同士の `%` 演算を実行した場合, the Sugu Evaluator shall 従来と同じ整数剰余結果を返す（例: `10 % 3` → `1`）
3. If ゼロで剰余演算が行われた場合, the Sugu Evaluator shall エラーメッセージを返す

---

### Requirement 2: delete 関数

**Objective:** 開発者として、マップからキーを削除できるようにしたい。これによりマップの動的な操作が可能になる。

#### Acceptance Criteria

1. When `delete(map, key)` が呼び出された場合, the Sugu Evaluator shall 指定されたキーをマップから削除し `true` を返す
2. When 存在しないキーが `delete` に渡された場合, the Sugu Evaluator shall エラーを発生させず `false` を返す
3. If `const` で宣言されたマップに対して `delete` が呼び出された場合, the Sugu Evaluator shall const 変数の変更を禁止するエラーを返す
4. If 第1引数がマップ以外の場合, the Sugu Evaluator shall 型エラーメッセージを返す

---

### Requirement 3: 文字列操作関数

**Objective:** 開発者として、文字列を柔軟に操作する組み込み関数がほしい。これにより実用的な文字列処理プログラムが書ける。

#### Acceptance Criteria

1. When `split(str, sep)` が呼び出された場合, the Sugu Evaluator shall 文字列を区切り文字で分割した配列を返す
2. When `join(arr, sep)` が呼び出された場合, the Sugu Evaluator shall 配列要素を区切り文字で結合した文字列を返す
3. When `trim(str)` が呼び出された場合, the Sugu Evaluator shall 前後の空白を除去した文字列を返す
4. When `replace(str, old, new)` が呼び出された場合, the Sugu Evaluator shall すべての一致箇所を置換した文字列を返す
5. When `substring(str, start, end)` が呼び出された場合, the Sugu Evaluator shall rune 単位で部分文字列を返す（マルチバイト文字対応）
6. When `indexOf(str, substr)` が呼び出された場合, the Sugu Evaluator shall 部分文字列の rune 単位の開始位置を返す（見つからない場合は `-1`）
7. When `toUpper(str)` が呼び出された場合, the Sugu Evaluator shall すべての文字を大文字に変換した文字列を返す
8. When `toLower(str)` が呼び出された場合, the Sugu Evaluator shall すべての文字を小文字に変換した文字列を返す
9. If 引数の型が不正な場合, the Sugu Evaluator shall 適切な型エラーメッセージを返す

---

### Requirement 4: 数学関数

**Objective:** 開発者として、基本的な数学計算を行う組み込み関数がほしい。これにより数値処理プログラムの表現力が向上する。

#### Acceptance Criteria

1. When `abs(x)` が呼び出された場合, the Sugu Evaluator shall 絶対値を返す
2. When `floor(x)` が呼び出された場合, the Sugu Evaluator shall 負の無限大方向への切り捨て結果を返す（`floor(-3.7)` → `-4`、`int(-3.7)` → `-3` との違いを維持）
3. When `ceil(x)` が呼び出された場合, the Sugu Evaluator shall 正の無限大方向への切り上げ結果を返す
4. When `round(x)` が呼び出された場合, the Sugu Evaluator shall 四捨五入した結果を返す
5. When `min(a, b, ...)` が可変長引数で呼び出された場合, the Sugu Evaluator shall 最小値を返す
6. When `max(a, b, ...)` が可変長引数で呼び出された場合, the Sugu Evaluator shall 最大値を返す
7. When `random()` が呼び出された場合, the Sugu Evaluator shall 0以上1未満の乱数を返す
8. When `sqrt(x)` が呼び出された場合, the Sugu Evaluator shall 平方根を返す
9. When `pow(x, y)` が呼び出された場合, the Sugu Evaluator shall x の y 乗を返す
10. If 数値以外の引数が渡された場合, the Sugu Evaluator shall 型エラーメッセージを返す

---

### Requirement 5: for-in ループ

**Objective:** 開発者として、配列やマップを簡潔にイテレートできる構文がほしい。これによりコレクション操作のコードが読みやすくなる。

#### Acceptance Criteria

1. When `for (item in arr) { ... }` が実行された場合, the Sugu Evaluator shall 配列の各要素を順にイテレートする
2. When `for (i, item in arr) { ... }` が実行された場合, the Sugu Evaluator shall インデックスと要素をペアでイテレートする
3. When `for (key in map) { ... }` が実行された場合, the Sugu Evaluator shall マップの各キーをイテレートする
4. When `for (key, value in map) { ... }` が実行された場合, the Sugu Evaluator shall キーと値をペアでイテレートする
5. The Sugu Evaluator shall for-in のイテレーション変数を `const` として扱い、ループ本体内での再代入を禁止する
6. While for-in ループの実行中に `break` が使用された場合, the Sugu Evaluator shall ループを即座に終了する
7. While for-in ループの実行中に `continue` が使用された場合, the Sugu Evaluator shall 次のイテレーションにスキップする
8. The Sugu Evaluator shall `IN` トークン、`ForInStatement` AST ノードを追加し、パイプライン全体（Token → Lexer → Parser → Evaluator）で for-in をサポートする

---

### Requirement 6: 仕様ドキュメントの更新

**Objective:** 開発者として、新機能に対応した最新の言語仕様ドキュメントを参照したい。これにより言語の使い方を正確に把握できる。

#### Acceptance Criteria

1. When 各 Step の実装が完了した場合, the Sugu Evaluator shall `docs/specification.md` に新機能の構文・動作を追記する
2. When 各 Step の実装が完了した場合, the Sugu Evaluator shall `docs/backlog.md` から完了した項目を削除する
3. The Sugu Evaluator shall 各 Step でテーブル駆動テスト（正常系・異常系・エッジケース・マルチバイト）を作成する
