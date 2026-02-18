# AI-DLC and Spec-Driven Development

Kiro-style Spec Driven Development implementation on AI-DLC (AI Development Life Cycle)

## Project Context

### Paths
- Steering: `.kiro/steering/`
- Specs: `.kiro/specs/`

### Steering vs Specification

**Steering** (`.kiro/steering/`) - Guide AI with project-wide rules and context
**Specs** (`.kiro/specs/`) - Formalize development process for individual features

### Active Specifications
- Check `.kiro/specs/` for active specifications
- Use `/kiro:spec-status [feature-name]` to check progress

## Development Guidelines
- Think in English, generate responses in Japanese. All Markdown content written to project files (e.g., requirements.md, design.md, tasks.md, research.md, validation reports) MUST be written in the target language configured for this specification (see spec.json.language).

## Minimal Workflow
- Phase 0 (optional): `/kiro:steering`, `/kiro:steering-custom`
- Phase 1 (Specification):
  - `/kiro:spec-init "description"`
  - `/kiro:spec-requirements {feature}`
  - `/kiro:validate-gap {feature}` (optional: for existing codebase)
  - `/kiro:spec-design {feature} [-y]`
  - `/kiro:validate-design {feature}` (optional: design review)
  - `/kiro:spec-tasks {feature} [-y]`
- Phase 2 (Implementation): `/kiro:spec-impl {feature} [tasks]`
  - `/kiro:validate-impl {feature}` (optional: after implementation)
- Progress check: `/kiro:spec-status {feature}` (use anytime)

## Development Rules
- 3-phase approval workflow: Requirements → Design → Tasks → Implementation
- Human review required each phase; use `-y` only for intentional fast-track
- Keep steering current and verify alignment with `/kiro:spec-status`
- Follow the user's instructions precisely, and within that scope act autonomously: gather the necessary context and complete the requested work end-to-end in this run, asking questions only when essential information is missing or the instructions are critically ambiguous.

## Steering Configuration
- Load entire `.kiro/steering/` as project memory
- Default files: `product.md`, `tech.md`, `structure.md`
- Custom files are supported (managed via `/kiro:steering-custom`)

---

# Sugu Language Project

## 概要
Sugu は「すぐ実行できる」JavaScript風のインタプリタ言語です。

## 技術スタック
- 実装言語: Go
- 種類: インタプリタ

## プロジェクト構成
```
sugu/
├── token/      # トークン定義
├── lexer/      # 字句解析器
├── ast/        # 抽象構文木
├── parser/     # 構文解析器
├── object/     # オブジェクトシステム
├── evaluator/  # 評価器
├── repl/       # 対話型実行環境
└── docs/       # 仕様書
```

## 開発ルール

### コーディング規約
- Go の標準的なスタイルに従う
- `go fmt` でフォーマット
- テストは `*_test.go` に記述

### テスト
- 各パッケージにテストを書く
- `go test ./...` で全テスト実行

### コミット
- 日本語でコミットメッセージを書く
- 機能単位で小さくコミット

## Sugu 言語の構文

### 変数宣言
```javascript
mut x = 10;      // 再代入可能
const PI = 3.14; // 再代入不可
```

### 関数定義
```javascript
func add(a, b) => {
    return a + b;
}
```

### コメント
```javascript
// 単一行コメント
//-- 複数行コメント --//
```

## 参照ドキュメント
- [docs/specification.md](docs/specification.md) - 言語仕様
- [docs/phase1.md](docs/phase1.md) - Phase 1 実装範囲
