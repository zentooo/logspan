# LogSpan

LogSpanは、Go言語向けの構造化ロギングライブラリです。HTTPリクエスト単位でログを集約する**コンテキストロガー**と、即座にログを出力する**ダイレクトロガー**の2つのモードを提供し、柔軟で使いやすいロギング機能を実現します。

## 🎯 主な特徴

- **デュアルモードロギング**: コンテキストベースとダイレクトの2つのロギングモード
- **構造化ログ出力**: JSON形式での一貫したログ出力
- **ミドルウェア機構**: ログ処理パイプラインのカスタマイズが可能
- **DataDog連携**: DataDog Standard Attributesに対応
- **HTTPミドルウェア**: Webアプリケーションでの自動ログ設定
- **並行処理安全**: goroutineセーフな実装
- **シンプルなAPI**: 直感的で使いやすいインターフェース

## 📦 インストール

```bash
go get github.com/zentooo/logspan
```

## 🚀 クイックスタート

### ダイレクトロガー（即時出力）

```go
package main

import "github.com/zentooo/logspan/pkg/logger"

func main() {
    // グローバルダイレクトロガーを使用
    logger.D.Infof("アプリケーションが開始されました")
    logger.D.Warnf("警告: %s", "設定ファイルが見つかりません")
    logger.D.Errorf("エラー: %v", err)
}
```

### コンテキストロガー（ログ集約）

```go
package main

import (
    "context"
    "github.com/zentooo/logspan/pkg/logger"
)

func main() {
    // コンテキストロガーの作成
    ctx := context.Background()
    contextLogger := logger.NewContextLogger()
    ctx = logger.WithLogger(ctx, contextLogger)

    // コンテキスト情報の追加
    logger.AddField(ctx, "request_id", "req-12345")
    logger.AddField(ctx, "user_id", "user-67890")

    // ログの記録
    logger.Infof(ctx, "リクエスト処理を開始")
    processRequest(ctx)
    logger.Infof(ctx, "リクエスト処理が完了")

    // 集約されたログの出力
    logger.FlushContext(ctx)
}

func processRequest(ctx context.Context) {
    logger.AddField(ctx, "step", "validation")
    logger.Debugf(ctx, "入力パラメータを検証中")
    logger.Infof(ctx, "入力検証が完了")
}
```

## 📖 詳細な使用方法

### 1. 初期化と設定

#### グローバル設定

```go
import "github.com/zentooo/logspan/pkg/logger"

func init() {
    config := logger.Config{
        MinLevel:         logger.DebugLevel,
        Output:           os.Stdout,
        EnableSourceInfo: true,
    }
    logger.Init(config)
}
```

#### 個別ロガーの設定

```go
// ダイレクトロガーの設定
directLogger := logger.NewDirectLogger()
directLogger.SetLevelFromString("WARN")
directLogger.SetOutput(logFile)

// コンテキストロガーの設定
contextLogger := logger.NewContextLogger()
contextLogger.SetLevel(logger.InfoLevel)
contextLogger.SetOutput(logFile)
```

### 2. ログレベル

LogSpanは5つのログレベルをサポートしています：

- `DEBUG`: 詳細なデバッグ情報
- `INFO`: 一般的な情報メッセージ
- `WARN`: 警告メッセージ
- `ERROR`: エラーメッセージ
- `CRITICAL`: 重大なエラーメッセージ

```go
logger.D.Debugf("デバッグ情報: %s", debugInfo)
logger.D.Infof("情報: %s", info)
logger.D.Warnf("警告: %s", warning)
logger.D.Errorf("エラー: %v", err)
logger.D.Criticalf("重大なエラー: %v", criticalErr)
```

### 3. コンテキスト操作

#### フィールドの追加

```go
// 単一フィールドの追加
logger.AddField(ctx, "user_id", "12345")
logger.AddField(ctx, "session_id", "session-abc")

// 複数フィールドの追加
logger.AddFields(ctx, map[string]interface{}{
    "request_id": "req-67890",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
})
```

#### コンテキストロガーのAPI

```go
// コンテキストを使用したログ記録
logger.Infof(ctx, "ユーザー %s がログインしました", userID)
logger.Debugf(ctx, "処理ステップ: %s", step)
logger.Errorf(ctx, "処理中にエラーが発生: %v", err)

// ログの出力（集約されたログを一度に出力）
logger.FlushContext(ctx)
```

### 4. HTTPミドルウェア

Webアプリケーションでの自動ログ設定：

```go
package main

import (
    "net/http"
    "github.com/zentooo/logspan/pkg/http_middleware"
    "github.com/zentooo/logspan/pkg/logger"
)

func main() {
    mux := http.NewServeMux()

    // ログミドルウェアの適用
    handler := http_middleware.LoggingMiddleware(mux)

    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // リクエスト情報は自動的に追加される
        logger.Infof(ctx, "ユーザー一覧を取得中")

        // 追加のコンテキスト情報
        logger.AddField(ctx, "query_params", r.URL.Query())

        // 処理...

        logger.Infof(ctx, "ユーザー一覧の取得が完了")
        // FlushContext は自動的に呼ばれる
    })

    http.ListenAndServe(":8080", handler)
}
```

### 5. ミドルウェア機構

ログ処理パイプラインをカスタマイズできます：

```go
// パスワードマスキングミドルウェア
func passwordMaskingMiddleware(entry *logger.LogEntry, next func(*logger.LogEntry)) {
    // メッセージ内のパスワードをマスク
    entry.Message = strings.ReplaceAll(entry.Message, "password=secret", "password=***")
    next(entry)
}

// ミドルウェアの登録
logger.AddMiddleware(passwordMaskingMiddleware)

// タグ追加ミドルウェア
logger.AddMiddleware(func(entry *logger.LogEntry, next func(*logger.LogEntry)) {
    entry.Tags = append(entry.Tags, "production")
    next(entry)
})
```

### 6. フォーマッター

#### JSONフォーマッター（デフォルト）

```go
contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewJSONFormatter())
```

#### DataDogフォーマッター

```go
import "github.com/zentooo/logspan/pkg/formatter"

contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewDataDogFormatter())
```

## 📋 ログ出力形式

### デフォルトJSON形式

```json
{
  "type": "request",
  "context": {
    "request_id": "req-12345",
    "user_id": "user-67890"
  },
  "runtime": {
    "severity": "INFO",
    "startTime": "2023-10-27T09:59:58.123456+09:00",
    "endTime": "2023-10-27T10:00:00.223456+09:00",
    "elapsed": 150,
    "lines": [
      {
        "timestamp": "2023-10-27T09:59:59.123456+09:00",
        "level": "INFO",
        "message": "リクエスト処理を開始",
        "fields": {
          "request_id": "req-12345",
          "user_id": "user-67890"
        },
        "tags": []
      }
    ]
  },
  "config": {
    "elapsedUnit": "ms"
  }
}
```

### DataDog Standard Attributes形式

```json
{
  "timestamp": "2023-10-27T10:00:00Z",
  "status": "info",
  "message": "リクエスト処理を開始",
  "logger": "logspan",
  "duration": 150,
  "request_id": "req-12345",
  "user_id": "user-67890"
}
```

## 🔧 設定オプション

### Config構造体

```go
type Config struct {
    // 最小ログレベル
    MinLevel LogLevel

    // 出力先
    Output io.Writer

    // ソースファイル情報の有効化
    EnableSourceInfo bool
}
```

### デフォルト設定

```go
config := logger.DefaultConfig()
// MinLevel: InfoLevel
// Output: os.Stdout
// EnableSourceInfo: false
```

## 📚 サンプルコード

詳細なサンプルコードは `examples/` ディレクトリにあります：

```bash
# ダイレクトロガーのサンプル
go run examples/direct_logger/main.go

# コンテキストロガーのサンプル
go run examples/context_logger/main.go

# HTTPミドルウェアのサンプル
go run examples/http_middleware_example.go
```

## 🧪 テスト

```bash
# 全テストの実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 詳細なテスト出力
go test -v ./...
```

## 🏗️ アーキテクチャ

### パッケージ構成

```
pkg/
├── logger/                 # メインロガーパッケージ
│   ├── logger.go          # コアインターフェースとAPI
│   ├── context_logger.go  # コンテキストロガー実装
│   ├── direct_logger.go   # ダイレクトロガー実装
│   ├── config.go          # 設定管理
│   ├── middleware.go      # ミドルウェア機構
│   └── context.go         # コンテキストヘルパー
├── formatter/             # フォーマッター
│   ├── interface.go       # フォーマッターインターフェース
│   ├── json_formatter.go  # JSONフォーマッター
│   └── datadog_formatter.go # DataDogフォーマッター
└── http_middleware/       # HTTPミドルウェア
    └── middleware.go
```

### 設計原則

1. **シンプルなAPI**: 直感的で使いやすいインターフェース
2. **柔軟性**: 様々な用途に対応できる設計
3. **拡張性**: ミドルウェアによるカスタマイズ
4. **パフォーマンス**: 効率的なログ処理
5. **並行処理安全**: goroutineセーフな実装

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 🔗 関連リンク

- [Go Documentation](https://pkg.go.dev/github.com/zentooo/logspan)
- [Examples](./examples/)
- [Design Document](./design.md)

## 📞 サポート

質問や問題がある場合は、[Issues](https://github.com/zentooo/logspan/issues) を作成してください。