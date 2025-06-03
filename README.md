# LogSpan

LogSpanは、Go言語向けの構造化ロギングライブラリです。HTTPリクエスト単位でログを集約する**コンテキストロガー**と、即座にログを出力する**ダイレクトロガー**の2つのモードを提供し、柔軟で使いやすいロギング機能を実現します。

## 🎯 主な特徴

- **デュアルモードロギング**: コンテキストベースとダイレクトの2つのロギングモード
- **構造化ログ出力**: JSON形式での一貫したログ出力
- **ソース情報取得**: 関数名、ファイル名、行番号の自動取得（デバッグ支援）
- **ミドルウェア機構**: ログ処理パイプラインのカスタマイズが可能
- **コンテキスト展開**: contextフィールドをトップレベルに展開するフォーマッター
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
    logger.AddContextValue(ctx, "request_id", "req-12345")
    logger.AddContextValue(ctx, "user_id", "user-67890")

    // ログの記録
    logger.Infof(ctx, "リクエスト処理を開始")
    processRequest(ctx)
    logger.Infof(ctx, "リクエスト処理が完了")

    // 集約されたログの出力
    logger.FlushContext(ctx)
}

func processRequest(ctx context.Context) {
    logger.AddContextValue(ctx, "step", "validation")
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
        EnableSourceInfo: true,  // ソース情報（関数名、ファイル名、行番号）を有効化
        PrettifyJSON:     true,  // 読みやすいJSON形式で出力
        MaxLogEntries:    1000,  // 1000エントリで自動フラッシュ
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

### 3. ソース情報機能

LogSpanは、デバッグやトラブルシューティングを支援するため、ログエントリにソースコード情報を自動的に追加する機能を提供します。

#### ソース情報の有効化

```go
// グローバル設定でソース情報を有効化
config := logger.Config{
    MinLevel:         logger.DebugLevel,
    EnableSourceInfo: true,  // ソース情報を有効化
    Output:           os.Stdout,
}
logger.Init(config)

// ログ出力時に自動的にソース情報が追加される
logger.D.Infof("アプリケーションが開始されました")
```

#### 出力されるソース情報

ソース情報が有効な場合、各ログエントリに以下の情報が自動的に追加されます：

- `funcname`: 関数名（パッケージ名を含む完全な関数名）
- `filename`: ファイル名（パスを除いたファイル名のみ）
- `fileline`: 行番号

```json
{
  "timestamp": "2023-10-27T09:59:59.123456+09:00",
  "level": "INFO",
  "message": "アプリケーションが開始されました",
  "funcname": "main.main",
  "filename": "main.go",
  "fileline": 15
}
```

#### 使用例

```go
package main

import (
    "context"
    "github.com/zentooo/logspan/pkg/logger"
)

func main() {
    // ソース情報を有効化
    config := logger.DefaultConfig()
    config.EnableSourceInfo = true
    logger.Init(config)

    // ダイレクトロガーでの使用
    logger.D.Infof("アプリケーション開始")  // main.main, main.go, 行番号が記録される

    // コンテキストロガーでの使用
    ctx := context.Background()
    contextLogger := logger.NewContextLogger()
    ctx = logger.WithLogger(ctx, contextLogger)

    processUser(ctx, "user123")
    logger.FlushContext(ctx)
}

func processUser(ctx context.Context, userID string) {
    logger.Infof(ctx, "ユーザー処理開始: %s", userID)  // main.processUser, main.go, 行番号が記録される

    validateUser(ctx, userID)
}

func validateUser(ctx context.Context, userID string) {
    logger.Debugf(ctx, "ユーザー検証中: %s", userID)  // main.validateUser, main.go, 行番号が記録される
}
```

#### パフォーマンスに関する注意

ソース情報の取得には `runtime.Caller()` を使用するため、わずかなパフォーマンスオーバーヘッドが発生します。本番環境では必要に応じて無効化することを検討してください：

```go
// 本番環境での設定例
config := logger.Config{
    MinLevel:         logger.InfoLevel,
    EnableSourceInfo: false,  // 本番環境では無効化
    Output:           logFile,
}
logger.Init(config)
```

#### デバッグ時の活用

ソース情報機能は、特に以下の場面で有用です：

- **デバッグ**: ログの出力元を素早く特定
- **トラブルシューティング**: 問題の発生箇所を正確に把握
- **コードレビュー**: ログの出力箇所を確認
- **開発環境**: 詳細な情報でデバッグ効率を向上

```go
// 開発環境での設定例
config := logger.Config{
    MinLevel:         logger.DebugLevel,
    EnableSourceInfo: true,   // 開発時は有効化
    PrettifyJSON:     true,   // 読みやすい形式で出力
    Output:           os.Stdout,
}
logger.Init(config)
```

### 4. コンテキスト操作

#### コンテキストロガーの設定

```go
// コンテキストロガーの作成と設定
ctx := context.Background()
contextLogger := logger.NewContextLogger()
ctx = logger.WithLogger(ctx, contextLogger)

// または、コンテキストから自動取得（存在しない場合は新規作成）
contextLogger := logger.FromContext(ctx)
```

#### コンテキストフィールドの追加

```go
// 単一フィールドの追加
logger.AddContextValue(ctx, "user_id", "12345")
logger.AddContextValue(ctx, "session_id", "session-abc")

// 複数フィールドの追加
logger.AddContextValues(ctx, map[string]interface{}{
    "request_id": "req-67890",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
})

// 直接ロガーインスタンスを使用
contextLogger := logger.FromContext(ctx)
contextLogger.AddContextValue("operation", "user_login")
contextLogger.AddContextValues(map[string]interface{}{
    "step": "validation",
    "attempt": 1,
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

### 5. ミドルウェア機構

ログ処理パイプラインをカスタマイズできます：

#### 基本的なミドルウェア

```go
// カスタムミドルウェアの作成
func customMiddleware(entry *logger.LogEntry, next func(*logger.LogEntry)) {
    // ログエントリの前処理
    entry.Message = "[CUSTOM] " + entry.Message

    // 次のミドルウェアまたは最終処理を呼び出し
    next(entry)
}

// ミドルウェアの登録
logger.AddMiddleware(customMiddleware)

// ミドルウェアの管理
logger.ClearMiddleware()                    // 全ミドルウェアをクリア
count := logger.GetMiddlewareCount()        // ミドルウェア数を取得
```

#### パスワードマスキングミドルウェア

LogSpanには、機密情報を自動的にマスクする組み込みミドルウェアが含まれています：

```go
// デフォルト設定でパスワードマスキングを有効化
passwordMasker := logger.NewPasswordMaskingMiddleware()
logger.AddMiddleware(passwordMasker.Middleware())

// カスタム設定でパスワードマスキング
passwordMasker := logger.NewPasswordMaskingMiddleware().
    WithMaskString("[REDACTED]").                           // マスク文字列をカスタマイズ
    WithPasswordKeys([]string{"password", "secret"}).       // マスク対象キーを設定
    AddPasswordKey("api_key").                              // 追加のキーを指定
    AddPasswordPattern(regexp.MustCompile(`token=\w+`))     // カスタム正規表現パターン

logger.AddMiddleware(passwordMasker.Middleware())

// 使用例
logger.D.Infof("User login: username=john password=secret123 token=abc123")
// 出力: "User login: username=john password=*** token=***"
```

##### デフォルトでマスクされるキーワード
- `password`, `passwd`, `pwd`, `pass`
- `secret`, `token`, `key`, `auth`
- `credential`, `credentials`, `api_key`
- `access_token`, `refresh_token`

##### サポートされるパターン
- `key=value` 形式: `password=secret` → `password=***`
- JSON形式: `"password":"secret"` → `"password":"***"`
- カスタム正規表現パターン

### 6. フォーマッター

#### JSONフォーマッター（デフォルト）

```go
contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewJSONFormatter())
```

#### ContextFlattenフォーマッター

```go
import "github.com/zentooo/logspan/pkg/formatter"

contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewContextFlattenFormatter())
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
        "message": "リクエスト処理を開始"
      }
    ]
  },
  "config": {
    "elapsedUnit": "ms"
  }
}
```

### ソース情報付きの出力形式

`EnableSourceInfo: true` の場合、各ログエントリにソース情報が追加されます：

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
        "funcname": "main.processRequest",
        "filename": "main.go",
        "fileline": 42
      }
    ]
  }
}
```

### Context Flatten形式

ContextFlattenフォーマッターを使用すると、contextフィールドがトップレベルに展開されます：

```json
{
  "request_id": "req-12345",
  "user_id": "user-67890",
  "type": "request",
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
        "funcname": "main.processRequest",
        "filename": "main.go",
        "fileline": 42
      }
    ]
  }
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

    // JSON出力の整形（インデント）を有効化
    PrettifyJSON bool

    // コンテキストロガーの最大エントリ数（0 = 制限なし）
    MaxLogEntries int
}
```

### デフォルト設定

```go
config := logger.DefaultConfig()
// MinLevel: InfoLevel
// Output: os.Stdout
// EnableSourceInfo: false
// PrettifyJSON: false
// MaxLogEntries: 1000
```

### カスタム設定例

```go
// 開発環境向け設定（整形されたJSON出力）
logger.Init(logger.Config{
    MinLevel:         logger.DebugLevel,
    Output:           os.Stdout,
    EnableSourceInfo: true,
    PrettifyJSON:     true,  // 読みやすい整形されたJSON
    MaxLogEntries:    500,   // 500エントリで自動フラッシュ
})

// 本番環境向け設定（コンパクトなJSON出力）
logger.Init(logger.Config{
    MinLevel:         logger.InfoLevel,
    Output:           logFile,
    EnableSourceInfo: false,
    PrettifyJSON:     false,  // コンパクトなJSON
    MaxLogEntries:    1000,   // 1000エントリで自動フラッシュ
})

// メモリ効率重視設定
logger.Init(logger.Config{
    MinLevel:      logger.InfoLevel,
    Output:        logFile,
    PrettifyJSON:  false,
    MaxLogEntries: 100,  // 頻繁な自動フラッシュでメモリ使用量を抑制
})

// 制限なし設定（手動フラッシュのみ）
logger.Init(logger.Config{
    MinLevel:      logger.InfoLevel,
    Output:        logFile,
    MaxLogEntries: 0,  // 自動フラッシュを無効化
})
```

### 設定の確認

```go
// ロガーが初期化されているかチェック
if logger.IsInitialized() {
    config := logger.GetConfig()
    fmt.Printf("Current log level: %s\n", config.MinLevel.String())
    fmt.Printf("Pretty JSON enabled: %t\n", config.PrettifyJSON)
    fmt.Printf("Max log entries: %d\n", config.MaxLogEntries)
}
```

## 🚀 メモリ最適化

### 自動フラッシュ機能

LogSpanは、メモリ使用量を制御するための自動フラッシュ機能を提供します：

#### 基本的な動作

```go
// 自動フラッシュの設定
logger.Init(logger.Config{
    MaxLogEntries: 100, // 100エントリで自動フラッシュ
})

ctx := context.Background()
contextLogger := logger.NewContextLogger()
ctx = logger.WithLogger(ctx, contextLogger)

logger.AddContextValue(ctx, "request_id", "req-123")

// 100エントリに達すると自動的にフラッシュされる
for i := 0; i < 250; i++ {
    logger.Infof(ctx, "Processing item %d", i)
}
// 結果: 2回の自動フラッシュ（100エントリ、200エントリ時点）
// 残り50エントリは手動フラッシュが必要

logger.FlushContext(ctx) // 残りのエントリを出力
```

#### 自動フラッシュの特徴

- **エントリカウント**: ログレベルフィルターを通過したエントリのみがカウントされます
- **バッチ処理**: 各自動フラッシュは独立したログバッチとして出力されます
- **コンテキスト保持**: コンテキストフィールドは自動フラッシュ後も保持されます
- **メモリ解放**: フラッシュ後、エントリは自動的にクリアされてメモリが解放されます

#### メモリ効率的な使用例

```go
// 大量ログ処理での設定例
logger.Init(logger.Config{
    MinLevel:      logger.InfoLevel,
    MaxLogEntries: 50,    // 小さなバッチサイズ
    PrettifyJSON:  false, // コンパクト出力
})

ctx := context.Background()
contextLogger := logger.NewContextLogger()
ctx = logger.WithLogger(ctx, contextLogger)

logger.AddContextValue(ctx, "batch_id", "batch-001")

// 大量データの処理
for i := 0; i < 10000; i++ {
    logger.Infof(ctx, "Processing record %d", i)

    if i%1000 == 0 {
        // 進捗をコンテキストに追加
        logger.AddContextValue(ctx, "progress", fmt.Sprintf("%d/10000", i))
    }
}
// 自動フラッシュにより、メモリ使用量は一定に保たれる

logger.FlushContext(ctx) // 最後の残りエントリを出力
```

#### 無効化オプション

```go
// 自動フラッシュを無効にする（従来の動作）
logger.Init(logger.Config{
    MaxLogEntries: 0, // 0 = 制限なし
})

// この場合、手動でFlushContext()を呼ぶまでエントリが蓄積される
```

## 📚 サンプルコード

詳細なサンプルコードは `examples/` ディレクトリにあります：

```bash
# ダイレクトロガーのサンプル
go run examples/direct_logger/main.go

# コンテキストロガーのサンプル
go run examples/context_logger/main.go

# 自動フラッシュ機能のサンプル
go run examples/auto_flush/main.go

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

# カバレッジレポートの生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### テストカバレッジ

LogSpanは高いテストカバレッジを維持しており、以下のような包括的なテストを提供しています：

#### 主要テストファイル
- **base_logger_test.go**: 共通ベースロガーの機能テスト
- **config_test.go**: 設定管理機能のテスト
- **context_test.go**: コンテキスト関連機能のテスト
- **context_logger_test.go**: コンテキストロガーの包括的テスト
- **direct_logger_test.go**: ダイレクトロガーの包括的テスト
- **middleware_manager_test.go**: ミドルウェア管理機能のテスト
- **formatter_utils_test.go**: フォーマット関連ユーティリティのテスト

#### テストの特徴
- **並行処理テスト**: goroutineセーフティの検証
- **エラーケーステスト**: 異常系の動作確認
- **カバレッジ最適化**: 重要な関数の100%カバレッジ
- **統合テスト**: 実際の使用パターンでの動作確認

## 🏗️ アーキテクチャ

### パッケージ構成

```
pkg/
├── logger/                          # メインロガーパッケージ
│   ├── logger.go                   # コアインターフェースとグローバルインスタンス
│   ├── base_logger.go              # 共通ベースロガー（共通機能の実装）
│   ├── context_logger.go           # コンテキストロガー実装
│   ├── direct_logger.go            # ダイレクトロガー実装
│   ├── middleware_manager.go       # グローバルミドルウェア管理
│   ├── formatter_utils.go          # フォーマット関連ユーティリティ
│   ├── config.go                   # 設定管理
│   ├── entry.go                    # ログエントリ構造
│   ├── middleware.go               # ミドルウェア機構
│   ├── context.go                  # コンテキストヘルパー
│   ├── level.go                    # ログレベル定義
│   └── password_masking_middleware.go # パスワードマスキング
├── formatter/                       # フォーマッター
│   ├── interface.go                # フォーマッターインターフェース
│   ├── json_formatter.go           # JSONフォーマッター
│   └── context_flatten_formatter.go # ContextFlattenフォーマッター
├── http_middleware/                 # HTTPミドルウェア
│   └── middleware.go               # HTTPリクエストロギング
└── examples/                        # 使用例
    ├── context_logger/             # コンテキストロガー例
    ├── direct_logger/              # ダイレクトロガー例
    ├── context_flatten_formatter/  # ContextFlattenフォーマッター例
    └── http_middleware_example.go  # HTTPミドルウェア例
```

### 設計原則

1. **シンプルなAPI**: 直感的で使いやすいインターフェース
2. **柔軟性**: 様々な用途に対応できる設計
3. **拡張性**: ミドルウェアによるカスタマイズ
4. **パフォーマンス**: 効率的なログ処理
5. **並行処理安全**: goroutineセーフな実装
6. **責任の分離**: 機能別にファイルを分離し、保守性を向上

### アーキテクチャの改善点

#### コード重複の削除
- **BaseLogger**: `DirectLogger`と`ContextLogger`の共通機能を`BaseLogger`に統合
- **共通メソッド**: `SetOutput`, `SetLevel`, `SetFormatter`などの重複実装を削除
- **一貫性**: mutex命名の統一とスレッドセーフティの向上

#### 責任の明確化
- **middleware_manager.go**: グローバルミドルウェア管理機能を独立したファイルに分離
- **formatter_utils.go**: フォーマット関連ユーティリティを独立したファイルに分離
- **logger.go**: コアインターフェースとグローバルインスタンスのみに集中

#### テストカバレッジの向上
- **新しいテストファイル**: `base_logger_test.go`, `config_test.go`, `context_test.go`などを追加
- **カバレッジ改善**: 未カバーだった関数（`IsInitialized`, `AddContextValues`など）をテスト対象に追加
- **並行テスト**: 並行処理の安全性を検証するテストを強化

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