package session

import (
	"github.com/gofrs/uuid"
	"github.com/liquanhui-99/lr"
)

type Manager struct {
	Store
	Propagator
}

// GetSession 获取Session信息
func (m Manager) GetSession(ctx *lr.Context) (Session, error) {
	sessId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

	return m.Get(ctx.Req.Context(), sessId)
}

// InitSession 初始化session信息
// 第一步：初始化uuid生成Session信息
// 第二步：id信息写入到响应中
// 第三步返回Session
func (m Manager) InitSession(ctx *lr.Context) (Session, error) {
	uid := uuid.Must(uuid.NewV4())
	sessId := uid.String()

	session, err := m.Generate(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}

	// 把session Id写入到响应中
	if err = m.Inject(sessId, ctx.Resp); err != nil {
		return nil, err
	}

	return session, nil
}

func (m Manager) Refresh(ctx *lr.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	return m.Store.Refresh(ctx.Req.Context(), sess.ID())
}

// RemoveSession 删除Session，用于登出操作
func (m Manager) RemoveSession(ctx *lr.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	// 删除响应中的session id
	if err = m.Propagator.Remove(ctx.Resp); err != nil {
		return err
	}

	return m.Store.Remove(ctx.Req.Context(), sess.ID())
}
