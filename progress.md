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
- [x] 3.1 共通化できるコードの抽出
- [x] 3.2 ファイル分割・統合の検討
- [x] 3.3 命名の統一案の作成
- [x] 3.4 構造改善案の作成

### 4. 実装
- [ ] 4.1 優先度の高い改善から実装
- [ ] 4.2 テストの実行と修正
- [ ] 4.3 ドキュメントの更新

## 進捗状況
- 開始日: 2024年12月19日
- 現在のフェーズ: 実装フェーズ開始
- 完了したフェーズ: 現状分析、問題点特定、改善案検討

## 3.1 共通化できるコードの抽出 - 完了

### 実装内容
1. **BaseLogger構造体の導入**
   - 新ファイル: `pkg/logger/base_logger.go`
   - 共通フィールド: `output`, `minLevel`, `formatter`, `mutex`
   - 共通メソッド: `SetOutput()`, `SetLevel()`, `SetFormatter()`, `SetLevelFromString()`, `isLevelEnabled()`
   - スレッドセーフなゲッター: `getOutput()`, `getFormatter()`

2. **DirectLoggerの修正**
   - BaseLoggerを埋め込み構造として使用
   - 重複メソッドを削除（SetOutput, SetLevel, SetFormatter, SetLevelFromString, isLevelEnabled）
   - mutexの名前を`mu`から`mutex`に統一

3. **ContextLoggerの修正**
   - BaseLoggerを埋め込み構造として使用
   - 重複メソッドを削除（SetOutput, SetLevel, SetFormatter, isLevelEnabled）
   - mutexの名前を`mu`から`mutex`に統一

4. **テストの追加**
   - 新ファイル: `pkg/logger/base_logger_test.go`
   - BaseLoggerの全機能をテスト
   - スレッドセーフティのテストも含む

### 効果
- **コード重複の削減**: 約50行の重複コードを削除
- **保守性の向上**: 共通機能の変更が一箇所で済む
- **一貫性の向上**: mutexの命名を統一
- **テストカバレッジの向上**: BaseLoggerの独立したテスト

### テスト結果
- 全てのテストが通過（既存 + 新規BaseLoggerテスト）
- 既存機能に影響なし

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
- **内部メソッド**: `isLevelEnabled`

## 3.2 ファイル分割・統合の検討 - 完了

### 現在のlogger.goの責任分析
1. **Loggerインターフェース定義** - 適切
2. **グローバルミドルウェア管理** - 分離すべき
3. **フォーマット処理** - 分離すべき
4. **グローバルロガーインスタンス** - 適切
5. **デフォルトフォーマッター作成** - 分離すべき

### 分割提案

#### 1. middleware_manager.go（新規作成）
**責任**: グローバルミドルウェア管理
**移動する機能**:
- `globalMiddlewareChain`, `middlewareMutex`, `middlewareOnce`
- `ensureMiddlewareChain()`
- `AddMiddleware()`, `ClearMiddleware()`, `GetMiddlewareCount()`
- `processWithGlobalMiddleware()`

**理由**: ミドルウェア管理は独立した責任として分離すべき

#### 2. formatter_utils.go（新規作成）
**責任**: フォーマット関連のユーティリティ
**移動する機能**:
- `createDefaultFormatter()`
- `formatLogOutput()`

**理由**: フォーマット処理は複数のロガーで使用される共通機能

#### 3. logger.go（残存）
**責任**: コアインターフェースとグローバルインスタンス
**残す機能**:
- `Logger`インターフェース定義
- `D`グローバルロガーインスタンス

### 統合検討

#### テストファイルの統合は不要
- 各テストファイルは適切なサイズ（200-450行）
- 機能別に分離されており、保守性が高い
- 統合するとテストの可読性が低下する

#### 小さなファイルの統合検討
- `entry.go`（48行）と`level.go`（95行）は独立性が高く、統合不要
- `context.go`（64行）は独立した責任を持ち、統合不要

### 分割後の構造
```
pkg/logger/
├── logger.go              # Loggerインターフェース、グローバルインスタンス
├── middleware_manager.go  # グローバルミドルウェア管理（新規）
├── formatter_utils.go     # フォーマット関連ユーティリティ（新規）
├── base_logger.go         # 共通ベースロガー
├── direct_logger.go       # 直接ログ出力
├── context_logger.go      # コンテキストベースログ
├── config.go              # グローバル設定管理
├── entry.go               # ログエントリ構造
├── level.go               # ログレベル定義
├── context.go             # コンテキストヘルパー
├── middleware.go          # ミドルウェアインターフェース
└── password_masking_middleware.go # パスワードマスキング
```

### 効果
- **責任の明確化**: 各ファイルが単一の責任を持つ
- **保守性の向上**: 関連機能の変更が局所化される
- **テスト容易性**: 機能別のテストが書きやすくなる
- **可読性の向上**: ファイルサイズが適切になる

### 実装優先度
1. **高**: middleware_manager.go の分離（グローバル状態の整理）
2. **中**: formatter_utils.go の分離（共通機能の整理）
3. **低**: その他の微調整

## 3.3 命名の統一案の作成 - 完了

### 現在の命名の問題点

#### 1. Mutex命名の不一致
- **BaseLogger**: `mutex sync.Mutex` ✅ 統一済み
- **グローバル変数**: `middlewareMutex`, `configMutex` - 接尾辞パターン

#### 2. 内部メソッドの命名
- **現在**: `isLevelEnabled()` (private) - 適切
- **公開関数**: `IsLevelEnabled()` (public) - 適切
- 両方とも必要で、適切に使い分けられている

#### 3. 変数名の一貫性
- **グローバル変数**: `globalMiddlewareChain`, `middlewareMutex`, `middlewareOnce`
- **設定関連**: `globalConfig`, `configMutex`, `initialized`

### 命名統一案

#### 1. Mutex命名規則の統一
**現在の状況**: 既に適切に統一されている
- **構造体内**: `mutex` (短縮形、適切)
- **グローバル変数**: `xxxMutex` (説明的、適切)

**理由**:
- 構造体内では文脈が明確なので短縮形が適切
- グローバル変数では説明的な名前が必要

#### 2. グローバル変数の命名パターン
**現在**: 一貫したパターンを使用
```go
// ミドルウェア関連
globalMiddlewareChain *MiddlewareChain
middlewareMutex       sync.RWMutex
middlewareOnce        sync.Once

// 設定関連
globalConfig Config
configMutex  sync.RWMutex
initialized  bool
```

**提案**: 現在の命名を維持（既に適切）

#### 3. 関数名の明確性
**現在の状況**: 適切に命名されている
- `processWithGlobalMiddleware` - 長いが明確
- `formatLogOutput` - 機能が明確
- `createDefaultFormatter` - 目的が明確

**提案**: 現在の命名を維持

#### 4. 定数の命名
**現在**: 適切に命名されている
```go
const (
    debugLevelString    = "DEBUG"
    infoLevelString     = "INFO"
    // ...
)
```

**提案**: 現在の命名を維持

### 改善が必要な箇所

#### 1. コメントの一貫性
**問題**: 一部のコメントで日本語と英語が混在
```go
// nilの出力先の場合は何もしない
if l.output == nil {
    return
}
```

**改善案**: 英語に統一
```go
// Do nothing if output is nil
if l.output == nil {
    return
}
```

#### 2. エラーメッセージの一貫性
**現在**: 英語で統一されている（適切）

### 命名規則ガイドライン

#### 1. 構造体フィールド
- **Mutex**: `mutex` (短縮形)
- **その他**: 完全な名前を使用

#### 2. グローバル変数
- **Mutex**: `xxxMutex` (説明的)
- **設定**: `globalXxx` または `xxx` + 説明

#### 3. 関数・メソッド
- **Public**: PascalCase
- **Private**: camelCase
- **説明的な名前**: 機能が明確になるよう命名

#### 4. 定数
- **文字列定数**: `xxxString` サフィックス
- **レベル定数**: `XxxLevel` サフィックス

### 実装不要
現在の命名は既に適切に統一されており、大きな変更は不要。
唯一の改善点は日本語コメントの英語化のみ。

## 3.4 構造改善案の作成 - 完了

### 全体的な構造改善提案

#### 1. アーキテクチャの改善
**現在の問題**:
- logger.goに複数の責任が集中
- フォーマット処理が分散
- グローバル状態の管理が複雑

**改善案**:
```
pkg/logger/
├── core/                   # コア機能（新規ディレクトリ）
│   ├── interface.go       # Loggerインターフェース
│   ├── base_logger.go     # 共通ベースロガー
│   └── global.go          # グローバルインスタンス
├── loggers/               # ロガー実装（新規ディレクトリ）
│   ├── direct_logger.go   # 直接ログ出力
│   └── context_logger.go  # コンテキストベースログ
├── middleware/            # ミドルウェア関連（新規ディレクトリ）
│   ├── manager.go         # グローバルミドルウェア管理
│   ├── interface.go       # ミドルウェアインターフェース
│   └── password_masking.go # パスワードマスキング
├── utils/                 # ユーティリティ（新規ディレクトリ）
│   ├── formatter.go       # フォーマット関連
│   ├── entry.go           # ログエントリ構造
│   └── level.go           # ログレベル定義
├── config.go              # グローバル設定管理
└── context.go             # コンテキストヘルパー
```

**メリット**:
- 責任の明確な分離
- 機能別のディレクトリ構成
- 新機能の追加が容易

**デメリット**:
- パッケージ構造が複雑になる
- インポートパスが長くなる
- 既存コードへの影響が大きい

#### 2. 現実的な改善案（推奨）
**フラットな構造を維持しつつ責任を分離**:
```
pkg/logger/
├── logger.go              # Loggerインターフェース、グローバルインスタンス
├── middleware_manager.go  # グローバルミドルウェア管理（新規）
├── formatter_utils.go     # フォーマット関連ユーティリティ（新規）
├── base_logger.go         # 共通ベースロガー（既存）
├── direct_logger.go       # 直接ログ出力
├── context_logger.go      # コンテキストベースログ
├── config.go              # グローバル設定管理
├── entry.go               # ログエントリ構造
├── level.go               # ログレベル定義
├── context.go             # コンテキストヘルパー
├── middleware.go          # ミドルウェアインターフェース
└── password_masking_middleware.go # パスワードマスキング
```

#### 3. 依存関係の改善
**現在の問題**:
- 循環依存のリスク
- グローバル状態への依存

**改善案**:
1. **設定の注入**: グローバル設定への直接アクセスを減らす
2. **インターフェースの活用**: 具象型への依存を減らす
3. **ファクトリーパターン**: ロガー作成の統一

#### 4. エラーハンドリングの改善
**現在の状況**: 基本的なエラーハンドリングは実装済み

**改善案**:
1. **エラー型の定義**: 独自エラー型の導入
2. **エラーコールバック**: エラー発生時のコールバック機能
3. **フォールバック機能**: 出力失敗時の代替手段

#### 5. パフォーマンスの改善
**現在の状況**: 基本的な最適化は実装済み

**改善案**:
1. **プールの活用**: LogEntryのオブジェクトプール
2. **バッファリング**: 出力のバッファリング最適化
3. **並行処理**: ミドルウェア処理の並行化（必要に応じて）

### 実装優先度

#### Phase 1: 高優先度の実装
1. **BaseLoggerの導入** ✅ 完了
2. **middleware_manager.goの分離** - 次のタスク
3. **formatter_utils.goの分離**
4. **日本語コメントの英語化**

#### Phase 2: 中優先度（次のリリースで実装）
1. **エラーハンドリングの改善**
2. **設定注入の改善**
3. **テストヘルパーの追加**

#### Phase 3: 低優先度（将来的に検討）
1. **ディレクトリ構造の再編成**
2. **パフォーマンス最適化**
3. **新機能の追加**

### 破壊的変更の回避
- 公開APIは変更しない
- 既存の使用方法は維持
- 内部実装のみ改善

### 改善効果の測定
1. **コード品質**: 循環複雑度、重複率
2. **保守性**: ファイルサイズ、責任の分離度
3. **テスト性**: テストカバレッジ、テスト実行時間
4. **パフォーマンス**: ログ出力速度、メモリ使用量

## 次のタスク（4.1 優先度の高い改善から実装）

### Phase 1: 高優先度の実装
1. **BaseLoggerの導入** ✅ 完了
2. **middleware_manager.goの分離** - 次のタスク
3. **formatter_utils.goの分離**
4. **日本語コメントの英語化**

### 実装方針
- 破壊的変更を避ける
- 既存のテストが通ることを確認
- 段階的に実装して影響を最小化