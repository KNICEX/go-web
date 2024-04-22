package session

import "context"

type Store interface {
	Get(ctx context.Context, id string) (Session, error)
	Set(ctx context.Context, sess Session) error
	Generate(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Refresh(ctx context.Context, id string) error
}

type StoreOption func(store Store)

func WithSessionBuilder(builder Builder) StoreOption {
	return func(store Store) {
		switch s := store.(type) {
		case *MemoStore:
			s.sessionBuilder = builder
		case *RedisStore:
			s.sessionBuilder = builder
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
