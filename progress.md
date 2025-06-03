# Context Flatten Formatter 作成タスク

## タスク概要
DataDogフォーマッターを置き換えて、`context`以下のキーをJSONのトップレベルに展開するだけのシンプルなフォーマッターを作成する。

## 要件分析
- 現在のLogOutput構造体: `{type, context, runtime, config}`
- 新フォーマッターの出力: `context`内のキーをトップレベルに展開し、他のフィールドもトップレベルに配置
- 例: `{"context": {"user_id": "123", "request_id": "req-456"}}` → `{"user_id": "123", "request_id": "req-456", "type": "request", "runtime": {...}, "config": {...}}`

## サブタスク

### 1. 新しいフォーマッターの設計
- [x] 1-1. フォーマッター名の決定（ContextFlattenFormatter）
- [x] 1-2. 出力構造の詳細設計
- [x] 1-3. キー衝突時の処理方針決定

#### 設計詳細
- **フォーマッター名**: `ContextFlattenFormatter`
- **出力構造**: `context`内のキー・値ペアをトップレベルに展開し、`type`, `runtime`, `config`もトップレベルに配置
- **キー衝突処理**: contextの値を優先（contextのキーが`type`, `runtime`, `config`と衝突した場合）
- **インデント対応**: JSONFormatterと同様にインデント機能を提供

### 2. フォーマッター実装
- [x] 2-1. `pkg/formatter/context_flatten_formatter.go`ファイル作成
- [x] 2-2. ContextFlattenFormatter構造体実装
- [x] 2-3. コンストラクタ関数実装（インデント対応）
- [x] 2-4. Format()メソッド実装
- [x] 2-5. context展開ロジック実装

### 3. テスト実装
- [x] 3-1. `pkg/formatter/context_flatten_formatter_test.go`ファイル作成
- [x] 3-2. 基本的なフォーマット機能のテスト
- [x] 3-3. インデント機能のテスト
- [x] 3-4. キー衝突時の動作テスト
- [x] 3-5. 空のcontextの処理テスト

### 4. 統合とドキュメント更新
- [x] 4-1. 既存のDataDogフォーマッター使用箇所の特定
- [x] 4-2. 使用例の作成
- [x] 4-3. README.mdの更新（DataDog関連記述の削除、ContextFlatten関連の追加）
- [x] 4-4. examples/README.mdの更新

### 5. 動作確認
- [x] 5-1. 全テストの実行確認
- [x] 5-2. 削除による影響がないことの確認

## 進め方
1. まずファイル削除を行う
2. 次にドキュメントを順次更新する
3. 最後に全体の動作確認を行う

## 進捗
- [x] プランニング完了
- [x] ファイル削除完了
- [x] ドキュメント更新完了
- [x] 動作確認完了

## 完了サマリー

### 削除されたファイル
- `pkg/formatter/datadog_formatter.go`
- `pkg/formatter/datadog_formatter_test.go`
- `examples/datadog_formatter/main.go`
- `examples/datadog_formatter/README.md`
- `examples/datadog_formatter/` ディレクトリ

### 更新されたドキュメント
- `.cursor/rules/api-usage.mdc`: DataDog関連記述を削除し、ContextFlattenFormatter記述に置き換え
- `.cursor/rules/log-format.mdc`: DataDog形式を削除し、ContextFlatten形式を追加
- `.cursor/rules/project-structure.mdc`: DataDog関連ファイル記述を削除し、ContextFlatten関連に置き換え
- `README.md`: DataDog関連記述を削除し、ContextFlattenFormatter記述に置き換え
- `examples/README.md`: DataDog関連記述を削除し、ContextFlattenFormatter記述に置き換え

### 動作確認結果
- 全テスト成功: `go test ./...` - PASS
- 全パッケージビルド成功: `go build ./...` - 成功
- 削除による影響なし

---

# DataDogフォーマッター削除タスク

## タスク概要
ContextFlattenFormatterの実装が完了したため、不要になったDataDogフォーマッターを削除し、関連するドキュメントを更新する。

## サブタスク

### 1. DataDogフォーマッター関連ファイルの削除
- [x] 1-1. `pkg/formatter/datadog_formatter.go`の削除
- [x] 1-2. `pkg/formatter/datadog_formatter_test.go`の削除
- [x] 1-3. `examples/datadog_formatter/`ディレクトリの削除

### 2. ドキュメント更新
- [x] 2-1. API Usage Guideの更新（DataDog関連記述の削除）
- [x] 2-2. Log Format Guideの更新（DataDog形式の削除、ContextFlatten形式の追加）
- [x] 2-3. Project Structure Guideの更新（DataDog関連ファイルの削除、ContextFlatten関連の追加）
- [x] 2-4. README.mdの更新（DataDog関連記述の削除、ContextFlatten関連の追加）
- [x] 2-5. examples/README.mdの更新

### 3. 動作確認
- [x] 3-1. 全テストの実行確認
- [x] 3-2. 削除による影響がないことの確認

## 進め方
1. まずファイル削除を行う
2. 次にドキュメントを順次更新する
3. 最後に全体の動作確認を行う

---

# ConfigInfo削除タスク

## タスク概要
ConfigInfoが不要になったため、全面的に削除し、関連するコードとテストを更新する。

## 要件分析
- ConfigInfoは現在`pkg/formatter/interface.go`で定義されている
- LogOutput構造体のConfigフィールドで使用されている
- 各フォーマッターのテストファイルで使用されている
- 実際の使用例でも使用されている

## サブタスク

### 1. ConfigInfo削除の影響調査
- [x] 1-1. ConfigInfoの使用箇所特定
- [x] 1-2. 削除による影響範囲の確認

### 2. コア実装の更新
- [x] 2-1. `pkg/formatter/interface.go`からConfigInfo定義とConfigフィールドを削除
- [x] 2-2. `pkg/logger/logger.go`でのConfigInfo使用箇所を削除

### 3. フォーマッター実装の更新
- [x] 3-1. `pkg/formatter/json_formatter.go`の更新（Configフィールド削除）
- [x] 3-2. `pkg/formatter/context_flatten_formatter.go`の更新（Configフィールド削除）

### 4. テストファイルの更新
- [x] 4-1. `pkg/formatter/json_formatter_test.go`の更新
- [x] 4-2. `pkg/formatter/context_flatten_formatter_test.go`の更新
- [x] 4-3. `pkg/logger/direct_logger_test.go`の更新

### 5. 使用例の更新
- [x] 5-1. `examples/context_flatten_formatter/main.go`の更新

### 6. ドキュメント更新
- [x] 6-1. API Usage Guideの更新（ConfigInfo関連記述の削除）
- [x] 6-2. Log Format Guideの更新（configフィールドの削除）

### 7. 動作確認
- [x] 7-1. 全テストの実行確認
- [x] 7-2. 削除による影響がないことの確認

## 進め方
1. まずコア実装（interface.go）を更新
2. 各フォーマッター実装を更新
3. テストファイルを更新
4. 使用例を更新
5. ドキュメントを更新
6. 最後に全体の動作確認を行う

## 進捗
- [x] プランニング完了
- [x] 実装完了

## 完了サマリー

### 削除されたコード
- `pkg/formatter/interface.go`: ConfigInfo構造体定義とLogOutput.Configフィールド
- `pkg/logger/logger.go`: ConfigInfo使用箇所
- `pkg/formatter/context_flatten_formatter.go`: configフィールド処理
- 全テストファイル: ConfigInfo使用箇所
- `examples/context_flatten_formatter/main.go`: ConfigInfo使用箇所

### 更新されたドキュメント
- `.cursor/rules/log-format.mdc`: configフィールドの記述を削除

### 動作確認結果
- 全テスト成功: `go test ./...` - PASS
- 全パッケージビルド成功: `go build ./...` - 成功
- 使用例正常動作: `go run examples/context_flatten_formatter/main.go` - 成功
- 削除による影響なし
