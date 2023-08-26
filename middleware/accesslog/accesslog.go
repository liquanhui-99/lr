package accesslog

import (
	"encoding/json"
	"github.com/liquanhui-99/lr"
)

type MiddleBuilder struct {
	// logFunc可以解决调用者使用不同log包的问题
	logFunc func(msg string)
}

func NewMiddleBuilder() *MiddleBuilder {
	return &MiddleBuilder{}
}

func (b *MiddleBuilder) LogFunc(fn func(string)) *MiddleBuilder {
	b.logFunc = fn
	return b
}

func (b *MiddleBuilder) Build() lr.Middleware {
	return func(next lr.HandleFunc) lr.HandleFunc {
		return func(ctx *lr.Context) {
			defer func() {
				// 记录请求信息，在defer中执行可以方式panic问题导致未知性
				lg := AccessLog{
					Host:       ctx.Req.Host,
					Root:       ctx.MatchedPath(),
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}

				bytes, _ := json.Marshal(&lg)
				b.logFunc(string(bytes))
			}()
			next(ctx)
		}
	}
}

type AccessLog struct {
	Host       string `json:"host,omitempty"`
	Root       string `json:"root,omitempty"`
	Path       string `json:"path,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
}
