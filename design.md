新しいGoロガーライブラリの設計書とパッケージ構成案を以下に示します。
これは、既存の `logone-go` のコンセプトを継承しつつ、指定された改善点を盛り込むことを目的としています。

## 新Goロガー 設計書

### 1. はじめに

本ドキュメントは、新たなGo言語向けロガーライブラリの設計を定義します。このライブラリは、`logone-go` の「HTTPリクエスト単位でログを単一JSONにまとめる」という優れたコンセプトを基盤としつつ、より柔軟な利用シーンへの対応、設定オプションの拡充、そして開発者にとってシンプルなインターフェースを提供することを目指します。

### 2. 基本コンセプト

* **デュアルモードロギング:**
    * **コンテキストロガー:** `context.Context` に紐付き、特定の処理単位（例: HTTPリクエスト）内のログエントリを内部的に集約し、処理完了時にまとめて出力します。
    * **ダイレクトロガー:** `context.Context` に依存せず、ログメソッド呼び出しの都度、ログを直接出力します。長時間稼働するバッチ処理などでの利用を想定します。
* **透過的なコンテキスト管理:** コンテキストロガーの初期化や終了処理（フラッシュ）は、可能な限りライブラリ側で自動化し、ユーザーが明示的に管理する手間を削減します（例: HTTPミドルウェア経由での自動セットアップ）。
* **シンプルなAPI:** 開発者が直感的に利用できるよう、ログ出力メソッドや設定インターフェースを簡潔に保ちます。
* **拡張性:** Logger Middleware機構を提供し、ログ出力パイプラインにカスタム処理（例: 機密情報のマスキング）を容易に組み込めるようにします。
* **標準への準拠オプション:** DataDog Standard Attributesなど、業界標準のログ形式に合わせた出力設定を可能にします（デフォルトではOFF）。

### 3. 主な機能とインターフェース

#### 3.1. 初期化

* **ライブラリ全体の初期化 (オプション):**
    ```go
    // logger.Init(config Config)
    // アプリケーション起動時に一度呼び出し、グローバルな設定（デフォルトログレベル、DataDog連携の有効化など）を行います。
    ```
    `Config` 構造体には、DataDog Standard Attributes を使用するかどうかのフラグを含めます。

* **コンテキストロガーの初期化:**
    * ユーザーによる明示的な `LogManager` の初期化は不要とします。
    * HTTPサーバーのミドルウェアなどを介して、リクエストコンテキストの開始時に自動的にロガーがセットアップされ、コンテキストに注入される仕組みを目指します。
    * ユーザーはコンテキストからロガーを取得して利用します。

#### 3.2. ログ出力インターフェース

* **コンテキストロガー:**
    * ログ出力メソッドは `context.Context` を第一引数に取ります。
    * 例: `logger.Infof(ctx context.Context, format string, args ...interface{})`
    * その他レベル (Debugf, Warnf, Errorf, Criticalf) も同様のシグネチャで提供します。

* **ダイレクトロガー:**
    * 専用のインターフェースまたはインスタンス経由でログ出力メソッドを呼び出します。
    * 例: `logger.D.Infof(format string, args ...interface{})` (ここで `logger.D` はダイレクトロガーのインスタンスまたは窓口)
    * その他レベルも同様。

    *初期化時のインターフェース共通化について:*
    両ロガーは、ログレベルごとのメソッド (`Infof`, `Errorf` 等) を持つという点で共通の振る舞いを持ちますが、コンテキストの有無で呼び出し方が異なります。上記のように、コンテキストロガーは `ctx` を渡し、ダイレクトロガーは別の窓口 (`D`) を使うことで区別します。

#### 3.3. ログコンテキスト操作 (コンテキストロガー向け)

* リクエスト処理の途中で、そのリクエストに紐づくログコンテキストに情報を追加できます。
    ```go
    // logger.AddField(ctx context.Context, key string, value interface{})
    // logger.AddFields(ctx context.Context, fields map[string]interface{})
    ```

#### 3.4. DataDog Standard Attributes 対応

* 初期設定 (`logger.Init`) またはコンテキストごとの設定で、ログ出力のキー名をDataDog Standard Attributesに準拠させるかを選択できます（デフォルトは非準拠）。
* 準拠させた場合、`http.request.method` や `http.status_code` などの標準的なキー名でログが出力されます。

#### 3.5. Logger Middleware

* ログエントリが最終的に出力される前に介在する処理層を定義できます。
    ```go
    /*
    type LogEntry struct { // ログエントリの構造 (仮)
        Timestamp time.Time
        Level     string
        Message   string
        Fields    map[string]interface{}
        // ...その他、ファイル名、行番号など
    }

    type Middleware func(entry *LogEntry, next func(*LogEntry))

    logger.AddMiddleware(yourMiddleware)
    */
    ```
* この仕組みを利用して、ユーザーは正規表現によるパスワードマスキングなどのカスタム処理をライブラリ本体とは独立して実装・適用できます。

### 4. パッケージ構成案

```
your-module-name/
├── logger/                     // メインパッケージ
│   ├── logger.go               // Loggerインターフェース定義、コアロジック、ユーザー向けAPI関数(Infof, D.Infofなど)
│   ├── context_logger.go       // コンテキストロガーの実装 (ログ集約機能)
│   ├── direct_logger.go        // ダイレクトロガーの実装 (即時出力機能)
│   ├── config.go               // 設定(Config)構造体、初期化処理(Init)
│   ├── entry.go                // LogEntry構造体定義
│   ├── middleware.go           // Middlewareインターフェースと処理基盤
│   ├── fields.go               // 共通フィールドキー定義、DataDog Standard Attributes定数
│   ├── context.go              // context.Contextへのロガー注入・取得ヘルパー
│   ├── formatter/              // 出力フォーマット関連
│   │   ├── interface.go        // Formatterインターフェース
│   │   ├── json_formatter.go   // 標準JSONフォーマッタ
│   │   └── datadog_formatter.go // DataDog準拠JSONフォーマッタ
│   └── http_middleware/        // (オプション) HTTPサーバー用ミドルウェア実装例
│       └── middleware.go
├── go.mod
└── README.md
// example/ (利用例)
//   main.go
```

### 5. ログ出力形式 (コンテキストロガーの例)

**デフォルト (logone-go 形式準拠):**
```json
{
  "type": "request",
  "context": {
    "REQUEST_ID": "xxxxxxx"
  },
  "runtime": {
    "severity": "INFO",
    "startTime": "2023-10-27T09:59:58.123456+09:00",
    "endTime": "2023-10-27T10:00:00.223456+09:00",
    "elapsed": 150,
    "lines": [
      {
        "severity": "DEBUG",
        "message": "invoked",
        "time": "2023-10-27T09:59:59.123456+09:00",
        "fileline": 19,
        "funcname": "main.main",
        "tags": [
          "critical", "debug_info"
        ],
        "attributes": {
            "key": "value",
            "process_id": 123
        }
      },
      {
        "severity": "INFO",
        "message": "Processing complete",
        "time": "2023-10-27T10:00:00.123456+09:00",
        "fileline": 25,
        "funcname": "main.processData",
        "tags": [
          "audit"
        ],
        "attributes": "Process successful"
      }
    ],
    "tags": {
      "critical": 1,
      "debug_info": 1,
      "audit": 1
    }
  },
  "config": {
    "elapsedUnit": "ms"
  }
}
```
**主なフィールド説明 (logone-go 準拠):**
* **`type`**: ログの種別 (例: "request")。
* **`context`**: リクエスト全体や処理単位で共通のコンテキスト情報。キーと値はユーザーが `logger.AddField` (または類似の機能) で設定します。`logone-go` の `LogContext` に相当します。
* **`runtime`**: リクエスト/処理単位の実行時情報。
    * `severity`: このコンテキスト内で記録された最も高いログレベル。
    * `startTime`: 処理開始時刻 (ISO8601形式)。
    * `endTime`: 処理終了時刻 (ISO8601形式)。
    * `elapsed`: 処理時間。単位は `config.elapsedUnit` に依存します。
    * `lines`: このコンテキスト内で記録された個々のログエントリの配列。
        * `severity`: 個々のログエントリのログレベル。
        * `message`: ログメッセージ。
        * `time`: ログ記録時刻 (ISO8601形式)。
        * `fileline`: ログが出力されたソースコードの行番号 (オプショナル)。
        * `funcname`: ログが出力された関数名 (オプショナル)。
        * `tags`: このログエントリに付与されたタグの配列 (オプショナル)。
        * `attributes`: このログエントリに付与された任意の追加情報 (オプショナル)。
    * `tags`: `lines` 内の全エントリの `tags` を集計し、各タグの出現回数を示したマップ。
* **`config`**: ロガーの設定情報。
    * `elapsedUnit`: `runtime.elapsed` の単位 (例: "ms", "ns")。

**DataDog Standard Attributes有効時 (一部抜粋):**
```json
{
  "timestamp": "2023-10-27T10:00:00Z",
  "status": "info",
  "message": "User logged in",
  // ... (変更なし)
}
```