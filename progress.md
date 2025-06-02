## サブタスク

### 1. 現在の examples 構造の確認
- [x] 1.1 examples ディレクトリの構造を確認
- [x] 1.2 既存の例（context_logger, direct_logger）の実装パターンを確認
- [x] 1.3 DataDog formatter の実装を確認

### 2. DataDog formatter example の作成
- [x] 2.1 examples/datadog_formatter/ ディレクトリを作成
- [x] 2.2 main.go ファイルを作成（context logger + DataDog formatter）
- [x] 2.3 README.md ファイルを作成（使用方法の説明）

### 3. 既存例の拡張
- [x] 3.1 context_logger/main.go にDataDog formatterの使用例を追加
- [x] 3.2 direct_logger/main.go にDataDog formatterの使用例を追加

### 4. ドキュメント更新
- [x] 4.1 examples/README.md を更新してDataDog formatter例を追加

## 進捗状況
- 開始日時: 2024年12月19日
- 完了日時: 2024年12月19日
- 状態: **完了**

## 完了したタスクの詳細
### 1.1 examples ディレクトリの構造確認
- examples/ ディレクトリには README.md, http_middleware_example.go, context_logger/, direct_logger/ が存在
- 各サブディレクトリには main.go ファイルが配置されている

### 1.2 既存例の実装パターン確認
- context_logger: コンテキストベースでログを蓄積し、FlushContextで一括出力
- direct_logger: 即座にログを出力、レベル設定のテスト機能付き

### 1.3 DataDog formatter実装確認
- pkg/formatter/datadog_formatter.go に実装済み
- DataDog Standard Attributes形式（timestamp, status, message, logger, duration）
- インデント付きフォーマットもサポート

### 2.1 examples/datadog_formatter/ ディレクトリ作成
- 新しいサンプル用のディレクトリを作成

### 2.2 main.go ファイル作成
- DataDog formatterを使用したcontext loggerの例を実装
- 標準JSON形式との比較も含む
- リクエスト処理のシミュレーション機能付き

### 2.3 README.md ファイル作成
- DataDog formatter使用方法の詳細説明
- 出力例とコードのポイントを記載
- DataDogでの活用方法も説明

### 3.1 context_logger/main.go 拡張
- 既存のcontext logger例にDataDog formatterセクションを追加
- formatter パッケージのimportを追加
- DataDog形式での出力例を実装

### 3.2 direct_logger/main.go 拡張
- 既存のdirect logger例にDataDog formatterセクションを追加
- 各ログレベルでのDataDog形式出力テストを実装

### 4.1 examples/README.md 更新
- DataDog formatter exampleの説明を追加
- ディレクトリ構造を更新
- フォーマッター説明セクションを追加

## 成果物
1. **新規作成**:
   - `examples/datadog_formatter/main.go` - DataDog formatter専用サンプル
   - `examples/datadog_formatter/README.md` - 詳細な使用方法説明

2. **既存ファイル拡張**:
   - `examples/context_logger/main.go` - DataDog formatter使用例追加
   - `examples/direct_logger/main.go` - DataDog formatter使用例追加
   - `examples/README.md` - DataDog formatter例の説明追加

## 実行確認
各サンプルは以下のコマンドで実行可能：
```bash
go run examples/datadog_formatter/main.go
go run examples/context_logger/main.go
go run examples/direct_logger/main.go
```
