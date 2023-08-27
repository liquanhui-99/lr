package errorHandler

import "github.com/liquanhui-99/lr"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: map[int][]byte{},
	}
}

func (b *MiddlewareBuilder) AddCode(code int, data []byte) *MiddlewareBuilder {
	b.resp[code] = data
	return b
}

func (b *MiddlewareBuilder) Build() lr.Middleware {
	return func(next lr.HandleFunc) lr.HandleFunc {
		return func(ctx *lr.Context) {
			next(ctx)
			respData, ok := b.resp[ctx.Status]
			if ok {
				// 串改结果
				ctx.RespData = respData
			}
		}
	}
}
