package session

import "context"

// Store 负责session的存储和创建
type Store interface {
	Get(ctx context.Context, id string) (Session, error)
	Set(ctx context.Context, sess Session) error
	Generate(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Refresh(ctx context.Context, id string) error
}

type StoreOption func(store Store)

func WithSessionCreator(builder Creator) StoreOption {
	return func(store Store) {
		switch s := store.(type) {
		case *MemoStore:
			s.sessionCreator = builder
		case *RedisStore:
			s.sessionCreator = builder
		}
	}
}

func WithSerializer(serializer Serializer) StoreOption {
	return func(store Store) {
		switch s := store.(type) {
		case *RedisStore:
			s.serializer = serializer
		case *MemoStore:
			panic("memo Store needn't serializer")
		}
	}
}
