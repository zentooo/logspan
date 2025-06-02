# Examples

このディレクトリには、logspanライブラリの使用例が含まれています。

## 実行方法

### Direct Logger Example
```bash
go run examples/direct_logger/main.go
```

Direct Loggerは即座にログを出力する機能を提供します。ログレベルによるフィルタリングや基本的なフォーマット機能をテストできます。DataDog formatterの使用例も含まれています。

### Context Logger Example
```bash
go run examples/context_logger/main.go
```

Context Loggerはリクエスト単位でログを集約し、JSON形式で出力する機能を提供します。コンテキストベースのログ管理をテストできます。DataDog formatterの使用例も含まれています。

### DataDog Formatter Example
```bash
go run examples/datadog_formatter/main.go
```

DataDog Standard Attributes形式でのログ出力に特化した例です。DataDogでの監視・分析に最適化されたログフォーマットの使用方法を学べます。

## ディレクトリ構造

```
examples/
├── README.md
├── direct_logger/
│   └── main.go          # Direct Logger の使用例（DataDog formatter含む）
├── context_logger/
│   └── main.go          # Context Logger の使用例（DataDog formatter含む）
└── datadog_formatter/
    ├── main.go          # DataDog Formatter 専用の使用例
    └── README.md        # DataDog Formatter の詳細説明
```

## フォーマッター

logspanライブラリは複数のログフォーマッターをサポートしています：

- **JSON Formatter（デフォルト）**: logone-go準拠の標準JSON形式
- **DataDog Formatter**: DataDog Standard Attributes形式

各exampleでは、これらのフォーマッターの使い分けを学ぶことができます。

## 注意事項

各exampleは独立したmain関数を持つため、個別のディレクトリに配置されています。これにより、`go test ./...` コマンドでのビルドエラーを回避しています。