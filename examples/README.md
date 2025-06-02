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

## ディレクトリ構造

```
examples/
├── README.md
├── direct_logger/
│   └── main.go          # Direct Logger の使用例
└── context_logger/
    └── main.go          # Context Logger の使用例
```

## 注意事項

各exampleは独立したmain関数を持つため、個別のディレクトリに配置されています。これにより、`go test ./...` コマンドでのビルドエラーを回避しています。