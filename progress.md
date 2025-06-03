# ログのtypeフィールド設定可能化タスク

## 概要
現在ハードコードされている`"request"`のtypeフィールドを設定可能にする機能を追加する。

## 全体の進め方
1. 現在の実装を分析し、typeフィールドがどこで設定されているかを特定
2. 設定可能にするための設計を決定（Config、Context、Logger単位など）
3. 必要なコード変更を実装
4. テストを追加・更新
5. ドキュメントを更新

## サブタスク

### 1. 現在の実装分析
- [x] 1-1. typeフィールドが設定されている箇所を特定 → pkg/logger/formatter_utils.go:42行目
- [x] 1-2. formatLogOutput関数の実装を詳細確認
- [x] 1-3. LogOutput構造体の使用箇所を確認

### 2. 設計決定
- [x] 2-1. 設定方法を決定（Config、Context、Logger単位のどれにするか） → Config単位で決定
- [x] 2-2. デフォルト値の扱いを決定 → "request"をデフォルトに
- [x] 2-3. 既存APIとの互換性を確保する方法を決定 → デフォルト値で互換性確保

### 3. 実装
- [x] 3-1. Config構造体にLogType設定を追加
- [x] 3-2. ContextLoggerにSetLogType機能を追加 → Config単位のため不要
- [x] 3-3. DirectLoggerにSetLogType機能を追加 → Config単位のため不要
- [x] 3-4. formatLogOutput関数を修正してtypeを動的に設定
- [x] 3-5. グローバル設定関数を追加 → 既存のInit関数で対応

### 4. テスト
- [x] 4-1. Config設定のテストを追加
- [x] 4-2. ContextLoggerのSetLogTypeテストを追加 → Config単位のため不要
- [x] 4-3. DirectLoggerのSetLogTypeテストを追加 → Config単位のため不要
- [x] 4-4. 既存テストが正常に動作することを確認

### 5. ドキュメント更新
- [ ] 5-1. API使用ガイドを更新
- [ ] 5-2. ログフォーマットガイドを更新
- [ ] 5-3. READMEを更新
- [ ] 5-4. 例を追加 → examples/log_type/main.go作成済み

## 進捗
- 開始日: 2024年12月19日
- 現在のステータス: 実装完了、ドキュメント更新中

## 実装内容
- Config構造体にLogTypeフィールドを追加（デフォルト: "request"）
- formatLogOutput関数でConfig.LogTypeを使用するように修正
- 空文字列の場合は"request"をデフォルトとして使用（後方互換性確保）
- テストを追加してカスタムLogTypeと空文字列の場合の動作を確認
- examples/log_type/main.goで使用例を作成
