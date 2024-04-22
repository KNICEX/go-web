package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/redis/go-redis/v9"
	"reflect"
	"time"
)

type RedisStore struct {
	cmd            redis.Cmdable
	exp            time.Duration
	sessionBuilder Builder
	serializer     Serializer
}

// Serializer 序列化Session接口
type Serializer interface {
	RegisterType(sess Session)
	Encode(Session) ([]byte, error)
	Decode([]byte) (Session, error)
}

type DefaultSerializer struct {
	sessType reflect.Type
}

func (d *DefaultSerializer) RegisterType(sess Session) {
	d.sessType = reflect.TypeOf(sess).Elem()
}

func (d *DefaultSerializer) Encode(sess Session) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(sess)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *DefaultSerializer) Decode(data []byte) (Session, error) {
	buf := bytes.NewBuffer(data)
	res := reflect.New(d.sessType)
	err := gob.NewDecoder(buf).DecodeValue(res)
	if err != nil {
		return nil, err
	}
	return res.Interface().(Session), nil
}

func NewRedisStore(cmd redis.Cmdable, expiration time.Duration, opts ...StoreOption) *RedisStore {
	res := &RedisStore{
		exp:            expiration,
		cmd:            cmd,
		sessionBuilder: DefaultBuilder,
		serializer:     &DefaultSerializer{},
	}
	for _, opt := range opts {
		opt(res)
	}
	res.serializer.RegisterType(res.sessionBuilder(res, ""))

	return res
}

func (r *RedisStore) Get(ctx context.Context, id string) (Session, error) {
	val, err := r.cmd.Get(ctx, id).Bytes()
	if err != nil {
		return nil, err
	}
	return r.serializer.Decode(val)
}

func (r *RedisStore) Set(ctx context.Context, sess Session) error {
	data, err := r.serializer.Encode(sess)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, sess.ID(), data, r.exp).Err()
}

func (r *RedisStore) Generate(ctx context.Context, id string) (Session, error) {
	sess := r.sessionBuilder(r, id)
	err := r.Set(ctx, sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (r *RedisStore) Remove(ctx context.Context, id string) error {
	return r.cmd.Del(ctx, id).Err()
}

func (r *RedisStore) Refresh(ctx context.Context, id string) error {
	return r.cmd.Expire(ctx, id, r.exp).Err()
}
