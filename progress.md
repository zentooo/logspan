# AddField関数の動作変更とリネーム

## タスク概要
現在のAddField関数は、LogEntryとcontextの両方にフィールドを追加しているが、contextのみに追加するように変更し、関数名をAddContextValueに変更する。

## 進め方
1. 現在の実装を詳細に調査
2. 変更が必要なファイルを特定
3. 段階的に変更を実施
4. テストを実行して動作確認

## サブタスク

### 1. 現在の実装調査
- [x] 1-1. AddField関数の実装詳細を確認
  - `pkg/logger/context.go`の`AddField`関数：contextからloggerを取得し、`logger.AddField`を呼び出す
  - `pkg/logger/context_logger.go`の`AddField`関数：`l.fields[key] = value`でcontextフィールドに追加
  - `addEntry`関数内で`l.fields`の内容を`entry.Fields`にコピーしている（74-76行目）
- [x] 1-2. LogEntryにフィールドが追加される箇所を特定
  - `pkg/logger/context_logger.go:81` - `addEntry`関数内でcontextフィールドをLogEntryにコピー
  - `pkg/logger/direct_logger.go:88` - DirectLoggerでLogEntry作成時に空のFieldsを初期化
  - ミドルウェア内でentry.Fields[key] = valueの形で直接追加される箇所が多数
- [x] 1-3. contextにフィールドが追加される箇所を特定
  - `pkg/logger/context_logger.go:94` - AddField関数内
  - `pkg/logger/context_logger.go:102` - AddFields関数内
- [x] 1-4. 影響範囲を調査（テスト、ドキュメント等）
  - **使用箇所**: examples/, pkg/http_middleware/, テストファイル多数
  - **ドキュメント**: README.md, design.md, .cursor/rules/api-usage.mdc
  - **主要な変更対象**:
    - 関数名変更: AddField → AddContextValue, AddFields → AddContextValues
    - 動作変更: LogEntryへのフィールドコピーを停止

### 2. 実装変更
- [x] 2-1. ContextLoggerのaddEntry関数を修正（contextフィールドをLogEntryにコピーしないように）
- [x] 2-2. AddField関数をAddContextValueに名前変更
- [x] 2-3. AddFields関数をAddContextValuesに名前変更（一貫性のため）

### 3. テスト修正
- [ ] 3-1. 関連するテストケースを修正
- [ ] 3-2. 新しい動作に合わせてテスト期待値を更新

### 4. ドキュメント更新
- [ ] 4-1. README.mdの更新
- [ ] 4-2. API使用ガイドの更新

### 5. 動作確認
- [ ] 5-1. 全テストの実行
- [ ] 5-2. 実際の動作確認
