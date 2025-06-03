# golangci-lint修正プランニング (更新版)

## 概要
test系とexamplesをignoreする設定に変更後、golangci-lintを実行して14個の問題を特定しました。これらの問題を分類し、優先度を付けて修正していきます。

**進捗状況**: 14個中9個完了、残り5個

## 問題の分類と修正計画

### 1. errcheck (4件) - 高優先度 ✅ 完了
エラーハンドリングが不適切な箇所（本体コードのみ）
- [x] 1.1 pkg/logger/context_logger.go:105 - fmt.Fprintf のエラーチェック
- [x] 1.2 pkg/logger/context_logger.go:109 - fmt.Fprintf のエラーチェック
- [x] 1.3 pkg/logger/direct_logger.go:96 - fmt.Fprintf のエラーチェック
- [x] 1.4 pkg/logger/direct_logger.go:100 - fmt.Fprintf のエラーチェック

### 2. unused (2件) - 高優先度 ✅ 完了
未使用の関数（デッドコード）
- [x] 2.1 pkg/logger/password_masking_middleware.go:106 - maskPasswordsInFields 関数削除
- [x] 2.2 pkg/logger/password_masking_middleware.go:129 - maskPasswordsInGenericMap 関数削除

### 3. goconst (5件) - 中優先度 ✅ 完了
定数化すべき文字列リテラル
- [x] 3.1 pkg/logger/level.go:23 - "DEBUG" を定数化
- [x] 3.2 pkg/logger/level.go:25 - "INFO" を定数化
- [x] 3.3 pkg/logger/level.go:27 - "WARN" を定数化
- [x] 3.4 pkg/logger/level.go:29 - "ERROR" を定数化
- [x] 3.5 pkg/logger/level.go:31 - "CRITICAL" を定数化

### 4. その他 (3件) - 低優先度 🔄 進行中
- [ ] 4.1 gochecknoinits: pkg/logger/logger.go:34 - init関数の使用見直し
- [ ] 4.2 gocritic: pkg/logger/password_masking_middleware.go:109 - strings.EqualFold の使用
- [ ] 4.3 goprintffuncname: pkg/logger/direct_logger.go:71 - printf関数の命名 (log → logf)

### 5. staticcheck (2件) - 低優先度 🔄 進行中
- [ ] 5.1 pkg/logger/context_logger.go:105 - 空のブランチ (SA9003)
- [ ] 5.2 pkg/logger/direct_logger.go:96 - 空のブランチ (SA9003)

## 実行方針
1. ✅ 高優先度の問題から順番に修正 → **完了**
2. ✅ 中優先度の問題を修正 → **完了**
3. 🔄 低優先度の問題を修正 → **進行中**
4. 各修正後にテストを実行して動作確認
5. 修正完了後に再度golangci-lintを実行して確認

## 修正の影響範囲
- **errcheck**: ログ出力時のエラーハンドリング改善 ✅ 完了
- **unused**: デッドコード削除によるコード品質向上 ✅ 完了
- **goconst**: ログレベル文字列の定数化による保守性向上 ✅ 完了
- **その他**: コードスタイルと命名規則の改善 🔄 進行中

## 修正完了状況
- ✅ **errcheck (4件)**: fmt.Fprintfのエラーハンドリングを適切に追加。フォールバック処理では意図的にエラーを無視するコメントを追加。
- ✅ **unused (2件)**: 未使用の関数maskPasswordsInFieldsとmaskPasswordsInGenericMapを削除。
- ✅ **goconst (5件)**: ログレベル文字列を定数として定義し、String()メソッドとParseLogLevel()関数で使用するように修正。

## 残り問題
現在5個の問題が残っています：
- gochecknoinits: 1件
- gocritic: 1件
- goprintffuncname: 1件
- staticcheck: 2件（空のブランチ）

## 次のステップ
低優先度の残り5件の問題を修正：
1. staticcheck問題（空のブランチ）の修正
2. gocritic問題（strings.EqualFold使用）の修正
3. goprintffuncname問題（関数名修正）の修正
4. gochecknoinits問題（init関数見直し）の修正
