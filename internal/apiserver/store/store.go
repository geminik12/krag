package store

import (
	"context"
	"sync"

	"github.com/geminik12/krag/pkg/db"
	"gorm.io/gorm"
)

var (
	once sync.Once
	S    IStore
)

type IStore interface {
	DB(ctx context.Context, wheres ...db.Where) *gorm.DB
	TX(ctx context.Context, fn func(ctx context.Context) error) error

	User() UserStore
}

type transactionKey struct {
}

type datastore struct {
	core *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) IStore {
	once.Do(func() {
		S = &datastore{db}
	})

	return S
}

func (store *datastore) DB(ctx context.Context, wheres ...db.Where) *gorm.DB {
	db := store.core

	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}

	for _, whr := range wheres {
		db = whr.Where(db)
	}

	return db
}

func (store *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return store.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// Users 返回一个实现了 UserStore 接口的实例.
func (store *datastore) User() UserStore {
	return newUserStore(store)
}
