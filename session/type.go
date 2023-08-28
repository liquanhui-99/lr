package session

import (
	"context"
	"net/http"
)

// Store 管理Session本身，声明周期的管理
type Store interface {
	// Generate id交给调用者取管理，过期时间交给业务
	Generate(ctx context.Context, id string) (Session, error)
	// Refresh 刷新Session
	Refresh(ctx context.Context, id string) error
	// Remove 删除Session
	Remove(ctx context.Context, id string) error
	// Get 获取Session信息
	Get(ctx context.Context, id string) (Session, error)
}

// Session session的核心接口
type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	ID() string
}

// Propagator 主要操作请求和响应
type Propagator interface {
	// Inject 把session注入到response中
	Inject(id string, writer http.ResponseWriter) error
	// Extract 从请求中取出session
	Extract(req *http.Request) (string, error)
	// Remove 从response中移除session
	Remove(write http.ResponseWriter) error
}
