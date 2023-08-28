package memory

import (
	"context"
	"errors"
	"github.com/liquanhui-99/lr/session"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("session: key找不到")
)

type Store struct {
	// 加锁保护缓存
	mu sync.RWMutex
	// 所有的Session都缓存本地cache中
	sessions *cache.Cache
	// 过期时间，统一设置
	expiration time.Duration
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions:   cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

// Generate 生成session，并添加到本地缓存
func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess := &Session{
		id:     id,
		values: map[string]any{},
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

// Refresh 刷新Session信息
func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mu.RLock()
	sess, ok := s.sessions.Get(id)
	s.mu.RUnlock()
	if !ok {
		return ErrKeyNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions.Set(id, sess, s.expiration)
	return nil
}

// Remove 删除缓存中的Session信息
func (s *Store) Remove(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions.Delete(id)
	return nil
}

// Get 获取缓存中的Session信息
func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, ErrKeyNotFound
	}

	return sess.(*Session), nil
}

type Session struct {
	values map[string]any
	mu     sync.RWMutex
	id     string
}

// Get 根据key从内存中获取Session信息
func (s *Session) Get(ctx context.Context, key string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.values[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return val, nil
}

// Set 设置Session到本地缓存
func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = val
	return nil
}

func (s *Session) ID() string {
	return s.id
}
