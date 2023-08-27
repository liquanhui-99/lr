package recover

import (
	"github.com/liquanhui-99/lr"
	"reflect"
)

type MiddlewareBuilder struct {
	// 响应码
	StatusCode int
	// 响应数据
	Data []byte
	// 记录panic发生时的日志，日志是可选的
	LogFunc func(ctx *lr.Context)
}

//func (b *MiddlewareBuilder) LogFunc(fn func(ctx *lr.Context)) *MiddlewareBuilder {
//	b.logFunc = fn
//	return b
//}

func (b *MiddlewareBuilder) Build() lr.Middleware {
	return func(next lr.HandleFunc) lr.HandleFunc {
		return func(ctx *lr.Context) {
			defer func() {
				if err := recover(); err != nil {
					// 重新赋值code和data数据即可
					ctx.Status = b.StatusCode
					ctx.RespData = b.Data
					// 判断logFunc是否已经赋值，如果赋值记录日志
					v := reflect.ValueOf(b.LogFunc)
					if v.Kind() == reflect.Func && !v.IsNil() {
						b.LogFunc(ctx)
					}
				}
			}()
			next(ctx)
		}
	}
}
