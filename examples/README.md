# Examples

このディレクトリには、logspanライブラリの使用例が含まれています。

## 実行方法

### Direct Logger Example
```bash
go run examples/direct_logger/main.go
```

Direct Loggerは即座にログを出力する機能を提供します。ログレベルによるフィルタリングや基本的なフォーマット機能をテストできます。

### Context Logger Example
```bash
go run examples/context_logger/main.go
```

Context Loggerはリクエスト単位でログを集約し、JSON形式で出力する機能を提供します。コンテキストベースのログ管理をテストできます。

### Auto Flush Example
```bash
go run examples/auto_flush/main.go
```

Auto Flush機能は、メモリ使用量を制御するためのエントリ数制限による自動フラッシュ機能を提供します。大量のログを扱う際のメモリ効率的な処理方法を学べます。

### Context Flatten Formatter Example
```bash
go run examples/context_flatten_formatter/main.go
```

Context Flatten Formatterは、contextフィールドをJSONのトップレベルに展開する形式でのログ出力に特化した例です。ログの可読性とアクセス性を向上させるフォーマットの使用方法を学べます。

### Middleware Example
```bash
go run examples/middleware/main.go
```

Middlewareの使用例では、パスワードマスキングミドルウェアの基本的な使用方法から、カスタム設定、複数ミドルウェアのチェーン処理まで学べます。機密情報の自動マスキングやログ処理のカスタマイズ方法を理解できます。

### Advanced Configuration Example
```bash
go run examples/advanced_config/main.go
```

Advanced Configurationでは、logspanライブラリの高度な設定オプションを学べます。ログレベルフィルタリング、カスタム出力先、自動フラッシュ設定、フォーマッター設定、Direct Loggerの詳細設定など、実用的な設定パターンを網羅しています。

### HTTP Middleware Example
```bash
go run examples/http_middleware_example.go
```

HTTP Middlewareは、Webアプリケーションでの自動ログ設定機能を提供します。HTTPリクエストの情報を自動的にコンテキストに追加し、リクエスト単位でのログ管理を実現します。

## ディレクトリ構造

```
examples/
├── README.md
├── direct_logger/
│   └── main.go          # Direct Logger の使用例
├── context_logger/
│   └── main.go          # Context Logger の使用例
├── auto_flush/
│   └── main.go          # Auto Flush 機能の使用例
├── context_flatten_formatter/
│   ├── main.go          # Context Flatten Formatter 専用の使用例
│   └── README.md        # Context Flatten Formatter の詳細説明
├── middleware/
│   └── main.go          # Middleware の使用例（パスワードマスキング等）
├── advanced_config/
│   └── main.go          # 高度な設定オプションの使用例
└── http_middleware_example.go  # HTTP Middleware の使用例
```

## 主な機能

### メモリ最適化
- **Auto Flush**: エントリ数制限による自動フラッシュでメモリ使用量を制御
- **バッチ処理**: 効率的なログ出力とメモリ管理
- **設定可能な制限**: アプリケーションの要件に応じた柔軟な設定

### ミドルウェアシステム
- **パスワードマスキング**: 機密情報の自動マスキング機能
- **カスタムミドルウェア**: 独自のログ処理ロジックの追加
- **ミドルウェアチェーン**: 複数のミドルウェアの組み合わせ
- **設定可能なマスキング**: カスタムキーワードとマスキング文字列

### フォーマッター

logspanライブラリは複数のログフォーマッターをサポートしています：

- **JSON Formatter（デフォルト）**: logone-go準拠の標準JSON形式
- **Context Flatten Formatter**: contextフィールドをトップレベルに展開する形式

各exampleでは、これらのフォーマッターの使い分けを学ぶことができます。

### 高度な設定オプション
- **ログレベルフィルタリング**: 動的なログレベル制御
- **カスタム出力先**: ファイル、標準出力、カスタムWriter
- **ソース情報**: ファイル名・行番号の自動追加
- **JSON整形**: 開発時の可読性向上オプション

## 注意事項

各exampleは独立したmain関数を持つため、個別のディレクトリに配置されています。これにより、`go test ./...` コマンドでのビルドエラーを回避しています。