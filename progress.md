# funcname と fileline 追加タスク

## 概要
現在のログ出力に `funcname`（関数名）と `fileline`（ファイル行番号）を追加する機能を実装する。

## 全体の進め方
1. 現在のログエントリ構造を分析
2. runtime.Callerを使用してソース情報を取得する仕組みを設計
3. LogEntry構造体を拡張
4. ソース情報取得機能を実装
5. 各ロガー（ContextLogger、DirectLogger）でソース情報を取得・設定
6. フォーマッターでの出力対応
7. テストケースの追加・更新
8. 設定オプションの追加（EnableSourceInfo）

## サブタスク

### 1. 現在の構造分析
- [x] 1-1. LogEntry構造体の確認
- [x] 1-2. フォーマッター構造の確認
- [x] 1-3. 既存のEnableSourceInfo設定の確認

### 2. ソース情報取得機能の設計・実装
- [x] 2-1. runtime.Callerを使用したソース情報取得関数の実装
- [x] 2-2. LogEntry構造体にfuncname、filename、filelineフィールドを追加
- [x] 2-3. フォーマッター側のLogEntry構造体も同様に拡張

### 3. ロガーでのソース情報取得・設定
- [x] 3-1. ContextLoggerでソース情報を取得・設定
- [x] 3-2. DirectLoggerでソース情報を取得・設定
- [x] 3-3. 適切なcaller skipレベルの調整

### 4. 設定オプションの活用
- [x] 4-1. EnableSourceInfo設定がtrueの時のみソース情報を取得するよう制御
- [x] 4-2. パフォーマンス考慮（runtime.Callerは重い処理のため）

### 5. テスト・検証
- [x] 5-1. 既存テストの更新
- [x] 5-2. 新しいテストケースの追加
- [x] 5-3. 実際の出力確認

### 6. ドキュメント更新
- [x] 6-1. README.mdにソース情報機能の説明を追加
- [x] 6-2. 設定オプションの説明を更新
- [x] 6-3. 出力形式の例を更新
- [x] 6-4. 使用例とベストプラクティスを追加

## 実装完了

✅ **すべてのサブタスクが完了しました！**

### 実装された機能
- `funcname`、`filename`、`fileline` フィールドをLogEntryに追加
- `EnableSourceInfo` 設定でON/OFF制御
- Context Logger、Direct Logger両方で対応
- 適切なcaller skipレベルの自動調整
- 包括的なテストケースの追加

### 出力例
```json
{
  "timestamp": "2025-06-03T11:24:55.153281+09:00",
  "level": "INFO",
  "message": "Direct logger test message",
  "funcname": "main.main",
  "filename": "test_source_info.go",
  "fileline": 19
}
```

### 使用方法
```go
// ソース情報を有効にする
config := logger.DefaultConfig()
config.EnableSourceInfo = true
logger.Init(config)

// 通常通りログを出力すると、ソース情報が自動的に含まれる
logger.D.Infof("test message")
```

## 技術的な考慮事項
- runtime.Caller()を使用してスタックトレースから関数名とファイル行を取得
- caller skipレベルの調整が必要（ログ関数→addEntry→実際の呼び出し元）
- EnableSourceInfo設定でON/OFF制御（パフォーマンス考慮）
- 既存のフォーマット構造との互換性維持

## 期待される出力形式
```json
{
  "timestamp": "2021-06-12T04:00:07.552964+09:00",
  "level": "DEBUG",
  "message": "invoked",
  "funcname": "main.main",
  "fileline": 19
}
```
