package session

import (
	"context"
	"github.com/fanjindong/go-cache"
	"time"
)

type MemoStore struct {
	store          cache.ICache
	exp            time.Duration
	sessionCreator Creator
}

func NewMemoStore(expiration time.Duration, opts ...StoreOption) *MemoStore {
	res := &MemoStore{
		store:          cache.NewMemCache(),
		exp:            expiration,
		sessionCreator: DefaultCreator,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (m *MemoStore) Get(ctx context.Context, id string) (Session, error) {
	sess, ok := m.store.Get(id)
	if !ok {
		return nil, ErrKeyNotFound
	}
	if s, ok := sess.(Session); ok {
		return s, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

func (m *MemoStore) Set(ctx context.Context, sess Session) error {
	ok := m.store.Set(sess.ID(), sess, cache.WithEx(m.exp))
	if !ok {
		return ErrSaveFailed
	}
	return nil
}

func (m *MemoStore) Generate(ctx context.Context, id string) (Session, error) {
	s := m.sessionCreator(m, id)
	ok := m.store.Set(id, s, cache.WithEx(m.exp))
	if !ok {
		return nil, ErrSaveFailed
	}
	return s, nil
}

func (m *MemoStore) Remove(ctx context.Context, id string) error {
	m.store.Del(id)
	return nil
}

func (m *MemoStore) Refresh(ctx context.Context, id string) error {
	ok := m.store.Expire(id, m.exp)
	if !ok {
		return ErrKeyNotFound
	}
	return nil
}
