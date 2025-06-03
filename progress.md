# loggerパッケージ整理タスク

## 目的
loggerパッケージ内のコードを分析し、整理・リファクタリングできる箇所を特定して改善する

## プランニング

### 1. 現状分析
- [x] 1.1 各ファイルの役割と責任を整理
- [x] 1.2 コードの重複箇所を特定
- [x] 1.3 命名規則の一貫性をチェック
- [x] 1.4 パッケージ構造の妥当性を検証

### 2. 問題点の特定
- [x] 2.1 重複コードの洗い出し
- [x] 2.2 責任の分離が不十分な箇所の特定
- [x] 2.3 命名の不一致や改善点の特定
- [x] 2.4 テストコードの整理が必要な箇所の特定

### 3. 改善案の検討
- [ ] 3.1 共通化できるコードの抽出
- [ ] 3.2 ファイル分割・統合の検討
- [ ] 3.3 命名の統一案の作成
- [ ] 3.4 構造改善案の作成

### 4. 実装
- [ ] 4.1 優先度の高い改善から実装
- [ ] 4.2 テストの実行と修正
- [ ] 4.3 ドキュメントの更新

## 進捗状況
- 開始日: 2024年12月19日
- 現在のフェーズ: 改善案の検討開始

## 1.1 各ファイルの役割と責任の整理結果

### 主要ファイルの役割
1. **logger.go** - メインAPI、グローバル関数、ミドルウェア管理、フォーマット処理
2. **config.go** - グローバル設定管理、初期化処理
3. **entry.go** - ログエントリ構造定義、ソース情報取得
4. **level.go** - ログレベル定義、レベル比較ユーティリティ
5. **middleware.go** - ミドルウェアインターフェース、チェーン管理
6. **context.go** - コンテキスト関連のヘルパー関数
7. **context_logger.go** - コンテキストベースのログ集約実装
8. **direct_logger.go** - 直接ログ出力実装
9. **password_masking_middleware.go** - パスワードマスキング機能

### 責任の分散状況
- **設定管理**: config.go + logger.go（createDefaultFormatter）
- **フォーマット処理**: logger.go（formatLogOutput）
- **ソース情報**: entry.go（getSourceInfo）
- **ミドルウェア**: middleware.go + logger.go（グローバル管理）
- **ログ出力**: context_logger.go + direct_logger.go

## 1.2 コードの重複箇所の特定結果

### 重複している機能・メソッド
1. **SetOutput メソッド**
   - `DirectLogger.SetOutput()` と `ContextLogger.SetOutput()`
   - 同じ実装パターン（mutex + フィールド設定）

2. **SetLevel メソッド**
   - `DirectLogger.SetLevel()` と `ContextLogger.SetLevel()`
   - 同じ実装パターン（mutex + フィールド設定）

3. **isLevelEnabled メソッド**
   - `DirectLogger.isLevelEnabled()` と `ContextLogger.isLevelEnabled()`
   - 両方とも `IsLevelEnabled(level, l.minLevel)` を呼び出し

4. **ログメソッドの実装パターン**
   - `Debugf`, `Infof`, `Warnf`, `Errorf`, `Criticalf`
   - 両ロガーで同じシグネチャ、`fmt.Sprintf(format, args...)` の使用

5. **createDefaultFormatter の使用**
   - `DirectLogger` と `ContextLogger` の初期化で同じ関数を呼び出し

## 1.3 命名規則の一貫性チェック結果

### 一貫している命名
- **ログメソッド**: `Debugf`, `Infof`, `Warnf`, `Errorf`, `Criticalf` - 全て統一
- **設定メソッド**: `SetOutput`, `SetLevel`, `SetFormatter` - 全て統一
- **レベル定数**: `DebugLevel`, `InfoLevel`, etc. - 全て統一

### 改善の余地がある命名
- **内部メソッド**: `isLevelEnabled` vs `IsLevelEnabled`（public/private の使い分け）
- **構造体フィールド**: 一部で略語使用（`mu` vs `mutex`）

## 1.4 パッケージ構造の妥当性検証結果

### 現在の構造
- **pkg/logger**: 16ファイル（10実装 + 6テスト）
- **pkg/formatter**: 6ファイル（4実装 + 2テスト）
- **pkg/http_middleware**: 3ファイル（2実装 + 1テスト）

### 構造の評価
- **適切な分離**: formatter と http_middleware は独立したパッケージとして適切
- **logger パッケージの肥大化**: 16ファイルは多めだが、機能的には妥当
- **テストカバレッジ**: 各主要機能にテストが存在

## 2.1 重複コードの洗い出し

### 高優先度の重複（共通化すべき）
1. **BaseLogger 構造体の導入候補**
   - `output io.Writer`
   - `minLevel LogLevel`
   - `formatter formatter.Formatter`
   - `mu sync.Mutex`
   - `SetOutput()`, `SetLevel()`, `SetFormatter()` メソッド

2. **共通ユーティリティメソッド**
   - `isLevelEnabled()` - 両ロガーで同じ実装
   - ソース情報取得ロジック - 異なるskipレベルだが基本パターンは同じ

### 中優先度の重複（検討が必要）
1. **ログメソッドのパターン**
   - `fmt.Sprintf(format, args...)` の使用
   - レベルチェック → エントリ作成 → 処理のフロー

## 2.2 責任の分離が不十分な箇所の特定

### 問題のある責任の混在
1. **logger.go の責任過多**
   - ミドルウェア管理（グローバル状態）
   - フォーマット処理（formatLogOutput）
   - デフォルトフォーマッター作成（createDefaultFormatter）
   - グローバルロガーインスタンス（D）

2. **config.go の設定更新処理**
   - グローバル設定管理
   - DirectLogger の直接操作（型アサーション）
   - フォーマッター作成の重複（createDefaultFormatterと同じロジック）

3. **フォーマット処理の分散**
   - logger.go: formatLogOutput関数
   - DirectLogger/ContextLogger: 各々でformatLogOutputを呼び出し
   - 変換ロジック（logger.LogEntry → formatter.LogEntry）の重複

## 2.3 命名の不一致や改善点の特定

### 不一致・改善が必要な命名
1. **構造体フィールドの略語使用**
   - `mu sync.Mutex` - 他の箇所では完全な名前を使用
   - 一貫性のため `mutex` に統一すべき

2. **内部メソッドの命名**
   - `isLevelEnabled` (private) vs `IsLevelEnabled` (public)
   - 機能は同じだが、privateメソッドは不要

3. **変数名の一貫性**
   - `middlewareMutex` vs `configMutex` vs `mu`
   - mutexの命名規則を統一すべき

4. **関数名の明確性**
   - `processWithGlobalMiddleware` - 長いが明確
   - `formatLogOutput` - 機能は明確だが配置場所が不適切

## 2.4 テストコードの整理が必要な箇所の特定

### テストの重複・改善点
1. **テストファイルのサイズ**
   - `direct_logger_test.go`: 436行 - 大きめ
   - `middleware_test.go`: 424行 - 大きめ
   - `context_logger_test.go`: 382行 - 大きめ

2. **共通テストパターンの重複**
   - `SetOutput(&buf)` の繰り返し使用
   - 両ロガーで同じテストパターン（SetOutput, SetLevel）
   - バッファ初期化とアサーションの重複

3. **テストヘルパーの不足**
   - 共通のセットアップ処理が各テストで重複
   - アサーション処理の重複
