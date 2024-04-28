package session

import (
	"github.com/KNICEX/go-web"
	"time"
)

const DefaultCtxSessionKey = "session-key"

// Manager session管理器，统筹session的生成、获取、删除等操作
type Manager struct {
	Propagator    Propagator
	Store         Store
	CtxSessionKey string
}

var DefaultManager = &Manager{
	Propagator:    NewCookiePropagator(),
	Store:         NewMemoStore(time.Minute * 30),
	CtxSessionKey: DefaultCtxSessionKey,
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	sess, ok := ctx.Get(m.CtxSessionKey)
	if ok {
		return sess.(Session), nil
	}

	sessId, err := m.Propagator.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	res, err := m.Store.Get(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}
	ctx.Set(m.CtxSessionKey, res)
	return res, nil
}

func (m *Manager) InitSession(ctx *web.Context, sessId string) (Session, error) {
	existId, err := m.Propagator.Extract(ctx.Req)
	if err == nil {
		_ = m.Store.Remove(ctx.Req.Context(), existId)
	}
	sess, err := m.Store.Generate(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}
	ctx.Set(m.CtxSessionKey, sess)
	// 注入http response
	err = m.Propagator.Inject(sess.ID(), ctx.Resp)
	return sess, err
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	ctx.Set(m.CtxSessionKey, nil)
	// response header中清除session信息
	return m.Propagator.Clean(ctx.Resp)
}

func (m *Manager) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Store.Refresh(ctx.Req.Context(), sess.ID())
}

func (m *Manager) SaveSession(ctx *web.Context, sess Session) error {
	err := m.Store.Set(ctx.Req.Context(), sess)
	if err != nil {
		return err
	}
	ctx.Set(m.CtxSessionKey, sess)
	return nil
}

func GetSession(ctx *web.Context) (Session, error) {
	return DefaultManager.GetSession(ctx)
}

func InitSession(ctx *web.Context, sessId string) (Session, error) {
	return DefaultManager.InitSession(ctx, sessId)
}

func RemoveSession(ctx *web.Context) error {
	return DefaultManager.RemoveSession(ctx)
}

func RefreshSession(ctx *web.Context) error {
	return DefaultManager.RefreshSession(ctx)
}

func SaveSession(ctx *web.Context, sess Session) error {
	return DefaultManager.SaveSession(ctx, sess)
}
