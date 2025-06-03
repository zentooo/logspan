# Logger Package Refactoring Plan

## 概要
loggerパッケージの整理とリファクタリングを行い、コードの重複を削減し、ログレベル比較ロジックを統一する。

## 主な改善点
1. ログレベル比較ロジックの統一と`level.go`への集約
2. 重複コードの削減
3. コードの可読性と保守性の向上
4. テストの整理と追加

## タスク一覧

### 1. ログレベル比較ロジックの統一
- [x] 1-1. `level.go`にログレベル比較メソッドを追加
  - [x] 1-1-1. `LogLevel`型に`GreaterThan(other LogLevel) bool`メソッドを追加
  - [x] 1-1-2. `LogLevel`型に`GreaterThanOrEqual(other LogLevel) bool`メソッドを追加
  - [x] 1-1-3. `LogLevel`型に`LessThan(other LogLevel) bool`メソッドを追加
  - [x] 1-1-4. `LogLevel`型に`LessThanOrEqual(other LogLevel) bool`メソッドを追加
- [x] 1-2. 文字列ベースの比較ロジックを`LogLevel`型ベースに変更
  - [x] 1-2-1. `logger.go`の`isHigherSeverity`関数を`LogLevel`型ベースに変更
  - [x] 1-2-2. `formatLogOutput`関数で最大severity計算を`LogLevel`型で行う
- [x] 1-3. テストの更新
  - [x] 1-3-1. `level.go`用の新しいテストファイル作成
  - [x] 1-3-2. 既存のテストを新しいAPIに合わせて更新

### 2. 重複コードの削減
- [x] 2-1. `isLevelEnabled`メソッドの統一
  - [x] 2-1-1. 共通の`isLevelEnabled`関数を`level.go`に移動
  - [x] 2-1-2. `ContextLogger`と`DirectLogger`で共通関数を使用
- [x] 2-2. フォーマッター初期化ロジックの統一
  - [x] 2-2-1. 共通のフォーマッター作成関数を作成
  - [x] 2-2-2. `ContextLogger`と`DirectLogger`で共通関数を使用

### 3. コード構造の改善
- [x] 3-1. `level.go`の機能拡張
  - [x] 3-1-1. ログレベル関連のユーティリティ関数を集約
  - [x] 3-1-2. ドキュメントコメントの充実
- [x] 3-2. `logger.go`の整理
  - [x] 3-2-1. 文字列ベースの比較ロジックを削除
  - [x] 3-2-2. 不要なコードの削除

### 4. テストの整理と追加
- [x] 4-1. 新しいテストファイルの作成
  - [x] 4-1-1. `level_test.go`の作成
- [x] 4-2. 既存テストの更新
  - [x] 4-2-1. `context_logger_test.go`の更新
  - [x] 4-2-2. `direct_logger_test.go`の更新
- [x] 4-3. テストカバレッジの確認

### 5. ドキュメントの更新
- [ ] 5-1. APIドキュメントの更新
- [ ] 5-2. 使用例の更新

## 実行方針
- 各サブタスクを順番に実行し、テストが通ることを確認してから次に進む
- 破壊的変更を避け、既存のAPIとの互換性を保つ
- リファクタリング後もすべてのテストが通ることを確認する

## 期待される効果
- コードの重複削減
- ログレベル比較の一貫性向上
- 保守性の向上
- パフォーマンスの向上（文字列比較からint比較へ）

## 完了した改善内容

### ✅ ログレベル比較ロジックの統一
- `LogLevel`型に比較メソッド（`GreaterThan`, `GreaterThanOrEqual`, `LessThan`, `LessThanOrEqual`）を追加
- `IsLevelEnabled`と`GetHigherLevel`のユーティリティ関数を追加
- 文字列ベースの`isHigherSeverity`関数を削除し、型安全な比較に変更
- `formatLogOutput`関数で`LogLevel`型ベースの最大severity計算に変更

### ✅ 重複コードの削減
- `ContextLogger`と`DirectLogger`の`isLevelEnabled`メソッドを共通の`IsLevelEnabled`関数を使用するように統一
- フォーマッター初期化ロジックを`createDefaultFormatter`関数に統一
- デッドロック問題を修正（`Init`関数内での循環ロック取得を回避）

### ✅ テストの整理と追加
- `level_test.go`を新規作成し、全ての新しいメソッドのテストを追加
- `context_logger_test.go`の`TestIsHigherSeverity`を新しいAPIに更新
- 全てのテストが正常に通ることを確認

### ✅ パフォーマンス向上
- 文字列比較からint比較への変更により、ログレベル比較のパフォーマンスが向上
- ハードコードされたレベルマップを削除し、型安全性を向上
