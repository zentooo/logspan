# DataDog Formatter Example

このサンプルは、logspanライブラリでDataDog Standard Attributes形式のログ出力を使用する方法を示しています。

## 概要

DataDog formatterを使用すると、ログ出力がDataDogの標準属性形式に変換され、DataDogでの分析やモニタリングが容易になります。

## 主な特徴

- **DataDog Standard Attributes**: `timestamp`, `status`, `message`, `logger`, `duration`などの標準フィールドを使用
- **カスタム属性**: コンテキストフィールドがDataDogのカスタム属性として追加される
- **ログレベル変換**: logspanのログレベル（INFO, WARN等）がDataDogのstatus（info, warn等）に自動変換
- **構造化ログ**: 複数のログエントリが構造化された形式で出力される

## 実行方法

```bash
cd examples/datadog_formatter
go run main.go
```

## 出力例

### DataDog形式の出力
```json
{
  "timestamp": "2023-10-27T10:00:00Z",
  "status": "info",
  "message": "Log context with multiple entries",
  "logger": "logspan",
  "duration": 150,
  "request_id": "req-12345",
  "user_id": "user-67890",
  "service": "user-service",
  "version": "1.2.3",
  "step": "validation",
  "lines": [
    {
      "timestamp": "2023-10-27T10:00:00.123Z",
      "status": "info",
      "message": "Starting request processing",
      "logger": "logspan"
    }
  ]
}
```

### 標準JSON形式との比較
標準のJSON形式では以下のような構造になります：
```json
{
  "type": "request",
  "context": {
    "request_id": "req-12345",
    "user_id": "user-67890"
  },
  "runtime": {
    "severity": "INFO",
    "startTime": "2023-10-27T10:00:00Z",
    "endTime": "2023-10-27T10:00:01Z",
    "elapsed": 150,
    "lines": [...]
  }
}
```

## コードのポイント

### 1. DataDog Formatterの設定
```go
contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewDataDogFormatterWithIndent("  "))
```

### 2. コンテキストフィールドの追加
```go
logger.AddContextValue(ctx, "request_id", "req-12345")
logger.AddContextValue(ctx, "service", "user-service")
```

### 3. ログの出力
```go
logger.Infof(ctx, "Processing started")
logger.FlushContext(ctx) // DataDog形式で出力
```

## DataDogでの活用

このフォーマットで出力されたログは、DataDogで以下のように活用できます：

- **フィルタリング**: `status:error` や `service:user-service` でフィルタ
- **ダッシュボード**: `duration` フィールドでパフォーマンス監視
- **アラート**: 特定の `status` や `request_id` に基づくアラート設定
- **ログ分析**: カスタム属性を使った詳細な分析

## 関連ファイル

- [main.go](./main.go) - サンプルコード
- [../../pkg/formatter/datadog_formatter.go](../../pkg/formatter/datadog_formatter.go) - DataDog formatter実装