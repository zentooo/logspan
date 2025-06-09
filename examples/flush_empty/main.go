package main

import (
	"context"
	"fmt"

	"github.com/zentooo/logspan/logger"
)

func main() {
	fmt.Println("=== FlushEmpty機能のデモンストレーション ===")

	// 1. デフォルト設定（FlushEmpty = true）
	fmt.Println("\n1. デフォルト設定（FlushEmpty = true）")
	logger.Init(logger.WithPrettifyJSON(true))

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	// コンテキスト情報を追加（HTTPリクエストのような情報）
	logger.AddContextValue(ctx, "request_id", "req-12345")
	logger.AddContextValue(ctx, "method", "GET")
	logger.AddContextValue(ctx, "path", "/api/users")
	logger.AddContextValue(ctx, "user_agent", "Mozilla/5.0")

	// ログエントリを追加しない状態でフラッシュ
	fmt.Println("ログエントリなしでフラッシュ:")
	logger.FlushContext(ctx)

	// 2. FlushEmpty無効化
	fmt.Println("\n2. FlushEmpty無効化（FlushEmpty = false）")
	logger.Init(
		logger.WithPrettifyJSON(true),
		logger.WithFlushEmpty(false),
	)

	ctx2 := context.Background()
	contextLogger2 := logger.NewContextLogger()
	ctx2 = logger.WithLogger(ctx2, contextLogger2)

	// 同じコンテキスト情報を追加
	logger.AddContextValue(ctx2, "request_id", "req-67890")
	logger.AddContextValue(ctx2, "method", "POST")
	logger.AddContextValue(ctx2, "path", "/api/users")

	// ログエントリを追加しない状態でフラッシュ
	fmt.Println("ログエントリなしでフラッシュ:")
	logger.FlushContext(ctx2)
	fmt.Println("（出力なし - FlushEmptyがfalseのため）")

	// 3. FlushEmpty有効化でログエントリもある場合
	fmt.Println("\n3. FlushEmpty有効化 + ログエントリあり")
	logger.Init(
		logger.WithPrettifyJSON(true),
		logger.WithFlushEmpty(true),
	)

	ctx3 := context.Background()
	contextLogger3 := logger.NewContextLogger()
	ctx3 = logger.WithLogger(ctx3, contextLogger3)

	// コンテキスト情報を追加
	logger.AddContextValue(ctx3, "request_id", "req-11111")
	logger.AddContextValue(ctx3, "method", "PUT")
	logger.AddContextValue(ctx3, "path", "/api/users/123")

	// ログエントリを追加
	logger.Infof(ctx3, "リクエスト処理開始")
	logger.Debugf(ctx3, "ユーザー情報を取得中")
	logger.Infof(ctx3, "リクエスト処理完了")

	fmt.Println("コンテキスト情報 + ログエントリありでフラッシュ:")
	logger.FlushContext(ctx3)

	fmt.Println("\n=== デモンストレーション完了 ===")
}
