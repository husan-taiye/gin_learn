package repository

import (
	"context"
	"gin_learn/webook/internal/domain/async"
)

type SMSAsyncReqRepository interface {
	Store(ctx context.Context, key string) error
	Find(ctx context.Context) (async.ReqAsync, error)
}

type DSMSAsyncReqRepository struct {
}
