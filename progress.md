# JSON出力Prettify機能追加タスク

## 概要
最終的なJSON出力をprettifyするためのconfigを作成し、examples以下で全てオンにする

## 現状分析
- 現在のConfig構造: MinLevel, Output, EnableSourceInfoの3つのフィールド
- Examples: context_logger, direct_logger, http_middleware_exampleの3つ
- JSON出力は現在compact形式（改行・インデントなし）

## タスク全体の進め方
1. Config構造にPrettifyJSONフィールドを追加
2. Formatter側でprettify対応を実装
3. 各exampleでprettify設定を有効化
4. 動作確認とテスト

## サブタスク一覧

### 1. Config構造の拡張
- [x] 1-1. pkg/logger/config.goにPrettifyJSONフィールドを追加
- [x] 1-2. DefaultConfig()でPrettifyJSONのデフォルト値を設定

### 2. Formatter側の対応
- [x] 2-1. pkg/formatter/json_formatter.goでprettify対応を実装
- [x] 2-2. pkg/formatter/datadog_formatter.goでprettify対応を実装
- [x] 2-3. Formatter interfaceの確認・必要に応じて拡張

### 3. Logger側の対応
- [ ] 3-1. ContextLoggerでconfig.PrettifyJSONを参照するよう修正
- [ ] 3-2. DirectLoggerでconfig.PrettifyJSONを参照するよう修正

### 4. Examples更新
- [ ] 4-1. examples/context_logger/main.goでprettify設定を有効化
- [ ] 4-2. examples/direct_logger/main.goでprettify設定を有効化
- [ ] 4-3. examples/http_middleware_example.goでprettify設定を有効化

### 5. 動作確認
- [ ] 5-1. 各exampleを実行してJSON出力が整形されることを確認
- [ ] 5-2. prettify無効時の動作も確認

## 進捗状況
- [ ] タスク開始