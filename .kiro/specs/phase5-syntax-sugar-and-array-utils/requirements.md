# Requirements Document

## Introduction

Phase 5 では、Sugu 言語の「書きやすさ」を向上させるための構文糖衣と配列ユーティリティを実装する。
具体的には、(1) インクリメント/デクリメント演算子と複合代入演算子、(2) 負のインデックスとスライス操作、(3) `contains`/`concat` 組み込み関数の 3 領域を対象とする。

これらの機能は JavaScript や Python などの主要言語で一般的であり、Sugu ユーザーが直感的にコードを書けるようにすることが目的である。

## Requirements

### Requirement 1: インクリメント・デクリメント演算子

**Objective:** Sugu ユーザーとして、`++` / `--` 演算子でカウンタ変数を簡潔に更新したい。`i = i + 1` のような冗長な記述を避けるためである。

#### Acceptance Criteria

1.1. When `x++` が式文として評価された場合、the Evaluator shall 変数 `x` の値を 1 増加させ、増加前の値を返す（後置インクリメント）。

1.2. When `x--` が式文として評価された場合、the Evaluator shall 変数 `x` の値を 1 減少させ、減少前の値を返す（後置デクリメント）。

1.3. If `++` または `--` の対象が `const` 変数である場合、the Evaluator shall 再代入エラーを返す。

1.4. If `++` または `--` の対象が数値型でない場合、the Evaluator shall 型エラーを返す。

1.5. The Parser shall `++` および `--` を後置演算子としてパースし、対応する AST ノードを生成する。

1.6. The Lexer shall `++` を `PLUS_PLUS` トークン、`--` を `MINUS_MINUS` トークンとして認識する。

### Requirement 2: 複合代入演算子

**Objective:** Sugu ユーザーとして、`+=`, `-=`, `*=`, `/=`, `%=` で変数を簡潔に更新したい。算術演算と代入を 1 つの式で記述するためである。

#### Acceptance Criteria

2.1. When `x += y` が評価された場合、the Evaluator shall `x` に `x + y` の結果を代入する。

2.2. When `x -= y` が評価された場合、the Evaluator shall `x` に `x - y` の結果を代入する。

2.3. When `x *= y` が評価された場合、the Evaluator shall `x` に `x * y` の結果を代入する。

2.4. When `x /= y` が評価された場合、the Evaluator shall `x` に `x / y` の結果を代入する。

2.5. When `x %= y` が評価された場合、the Evaluator shall `x` に `x % y` の結果を代入する。

2.6. If 複合代入の対象が `const` 変数である場合、the Evaluator shall 再代入エラーを返す。

2.7. If `/=` または `%=` の右辺が 0 である場合、the Evaluator shall ゼロ除算エラーを返す。

2.8. When `arr[i] += y` のようにインデックス式に対して複合代入が使われた場合、the Evaluator shall 対応する要素を更新する。

2.9. The Lexer shall `+=`, `-=`, `*=`, `/=`, `%=` をそれぞれ対応するトークンとして認識する。

2.10. The Parser shall 複合代入演算子を代入式としてパースし、対応する AST ノードを生成する。

### Requirement 3: 負のインデックス

**Objective:** Sugu ユーザーとして、`arr[-1]` で配列や文字列の末尾要素にアクセスしたい。末尾からの相対位置指定により、`arr[len(arr) - 1]` のような冗長な記述を避けるためである。

#### Acceptance Criteria

3.1. When 配列に負のインデックス `n`（n < 0）でアクセスした場合、the Evaluator shall `len(arr) + n` の位置にある要素を返す。

3.2. When 文字列に負のインデックス `n`（n < 0）でアクセスした場合、the Evaluator shall rune 単位で `len(str) + n` の位置にある文字を返す。

3.3. If 負のインデックスの絶対値が配列の長さを超える場合、the Evaluator shall `null` を返す。

3.4. If 負のインデックスの絶対値が文字列の長さ（rune 数）を超える場合、the Evaluator shall `null` を返す。

3.5. When 配列の負のインデックスに代入した場合、the Evaluator shall `len(arr) + n` の位置の要素を更新する。

3.6. If 代入時に負のインデックスの絶対値が配列の長さを超える場合、the Evaluator shall 範囲外エラーを返す。

### Requirement 4: スライス操作

**Objective:** Sugu ユーザーとして、`arr[1:3]` で配列や文字列の部分要素を取得したい。簡潔な構文で部分抽出を行うためである。

#### Acceptance Criteria

4.1. When `arr[start:end]` が評価された場合、the Evaluator shall インデックス `start` から `end - 1` までの要素を含む新しい配列を返す。

4.2. When `str[start:end]` が評価された場合、the Evaluator shall rune 単位でインデックス `start` から `end - 1` までの文字を含む新しい文字列を返す。

4.3. When `arr[:end]` のように start が省略された場合、the Evaluator shall インデックス `0` から `end - 1` までの要素を返す。

4.4. When `arr[start:]` のように end が省略された場合、the Evaluator shall インデックス `start` から末尾までの要素を返す。

4.5. When `arr[:]` のように両方が省略された場合、the Evaluator shall 全要素のコピーを返す。

4.6. When スライスに負のインデックスが使われた場合、the Evaluator shall 長さを基準に正のインデックスに変換して処理する。

4.7. If スライスの範囲が配列・文字列の長さを超える場合、the Evaluator shall 利用可能な範囲に自動クランプする（エラーにしない）。

4.8. The Evaluator shall スライス操作で元の配列・文字列を変更せず、常に新しいオブジェクトを返す。

4.9. The Parser shall `[` と `]` の間に `:` が含まれる式をスライス式としてパースし、対応する AST ノードを生成する。

### Requirement 5: contains 関数

**Objective:** Sugu ユーザーとして、配列や文字列に特定の要素・部分文字列が含まれるかを簡潔に判定したい。`indexOf` の戻り値を比較する間接的な方法ではなく、直接的に真偽値を得るためである。

#### Acceptance Criteria

5.1. When `contains(arr, value)` が評価された場合、the Evaluator shall 配列内に `value` と等しい要素があれば `true` を返す。

5.2. When `contains(arr, value)` が評価され配列内に該当要素がない場合、the Evaluator shall `false` を返す。

5.3. When `contains(str, substr)` が評価された場合、the Evaluator shall 文字列内に `substr` が含まれていれば `true` を返す。

5.4. When `contains(str, substr)` が評価され部分文字列が見つからない場合、the Evaluator shall `false` を返す。

5.5. If `contains` の第 1 引数が配列でも文字列でもない場合、the Evaluator shall 型エラーを返す。

5.6. If `contains` の引数が 2 つでない場合、the Evaluator shall 引数エラーを返す。

### Requirement 6: concat 関数

**Objective:** Sugu ユーザーとして、複数の配列を結合したい。`push` を繰り返す冗長な方法ではなく、1 回の呼び出しで結合するためである。

#### Acceptance Criteria

6.1. When `concat(arr1, arr2)` が評価された場合、the Evaluator shall 2 つの配列を結合した新しい配列を返す。

6.2. When `concat(arr1, arr2, arr3, ...)` のように 3 つ以上の配列が渡された場合、the Evaluator shall すべてを順に結合した新しい配列を返す。

6.3. The Evaluator shall `concat` の戻り値として常に新しい配列を返し、元の配列を変更しない。

6.4. If `concat` の引数のいずれかが配列でない場合、the Evaluator shall 型エラーを返す。

6.5. If `concat` の引数が 2 つ未満の場合、the Evaluator shall 引数エラーを返す。

6.6. When 空配列が `concat` に含まれる場合、the Evaluator shall 空配列を無視して正しく結合する。

### Requirement 7: ドキュメント更新

**Objective:** Sugu の言語仕様ドキュメントが常に最新の状態であるようにしたい。ユーザーが正確なリファレンスを参照できるようにするためである。

#### Acceptance Criteria

7.1. When Phase 5 の各機能が実装された場合、`docs/specification.md` shall 新しい演算子・構文・組み込み関数の仕様を含む。

7.2. When Phase 5 の各機能が実装された場合、`docs/index.md` shall 新機能のクイックリファレンスを含む。

7.3. When Phase 5 の各機能が実装された場合、`docs/backlog.md` shall 該当する課題が削除された状態である。
