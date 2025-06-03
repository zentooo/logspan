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
└── http_middleware_example.go  # HTTP Middleware の使用例
```

## 主な機能

### メモリ最適化
- **Auto Flush**: エントリ数制限による自動フラッシュでメモリ使用量を制御
- **バッチ処理**: 効率的なログ出力とメモリ管理
- **設定可能な制限**: アプリケーションの要件に応じた柔軟な設定

### フォーマッター

logspanライブラリは複数のログフォーマッターをサポートしています：

- **JSON Formatter（デフォルト）**: logone-go準拠の標準JSON形式
- **Context Flatten Formatter**: contextフィールドをトップレベルに展開する形式

各exampleでは、これらのフォーマッターの使い分けを学ぶことができます。

## 注意事項

各exampleは独立したmain関数を持つため、個別のディレクトリに配置されています。これにより、`go test ./...` コマンドでのビルドエラーを回避しています。