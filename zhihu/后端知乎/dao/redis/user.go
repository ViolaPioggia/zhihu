package redis

import (
	"context"
	"fmt"
	"main/app/global"
	"time"
)

func GetPassword(ctx context.Context, username string) (string, error) {
	GetKey := global.Rdb.Get(ctx, fmt.Sprintf("%s:password", username))
	if GetKey.Err() != nil {
		return "", GetKey.Err()
	}
	return GetKey.Val(), nil
}

func Set(ctx context.Context, username string, context string, expiration time.Duration) error {
	SetKV := global.Rdb.Set(ctx, username, context, expiration)
	return SetKV.Err()
}
