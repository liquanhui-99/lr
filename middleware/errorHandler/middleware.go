package errorHandler

import (
	"github.com/liquanhui-99/lr"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

type MiddlewareBuilder struct {
	// 返回内容，key为http标准响应码，value为模版引擎的
	resp map[int]lr.TemplateEngine
	// logFunc 日志的打印方式
	logFunc func(s any)
}

func NewMiddlewareBuilder(l ...func(any)) *MiddlewareBuilder {
	var fn = func(s any) {
		log.Fatal(s)
	}

	if len(l) > 0 {
		val := reflect.ValueOf(l[0])
		if val.Kind() != reflect.Func && val.IsNil() {
			fn = l[0]
		}
	}

	return &MiddlewareBuilder{
		resp:    map[int]lr.TemplateEngine{},
		logFunc: fn,
	}
}

func (b *MiddlewareBuilder) AddCode(code int) *MiddlewareBuilder {
	switch code {
	case http.StatusNotFound:
		notFoundEngine, err := template.ParseGlob("../../testdata/404.gohtml")
		if err != nil {
			b.logFunc(err.Error())
			return b
		}

		b.resp[code] = &lr.GoTemplateEngine{
			T: notFoundEngine,
		}
	case http.StatusInternalServerError:
		internalErrEngine, err := template.ParseGlob("../../testdata/500.gohtml")
		if err != nil {
			b.logFunc(err.Error())
			return b
		}

		b.resp[code] = &lr.GoTemplateEngine{
			T: internalErrEngine,
		}
	}

	return b
}

func (b *MiddlewareBuilder) Build() lr.Middleware {
	return func(next lr.HandleFunc) lr.HandleFunc {
		return func(ctx *lr.Context) {
			next(ctx)
			engine, ok := b.resp[ctx.Status]
			if ok {
				// 串改结果
				ctx.TplEngine = engine
				switch ctx.Status {
				case http.StatusNotFound:
					if err := ctx.Render("404.gohtml", nil); err != nil {
						b.logFunc(err.Error())
						return
					}
				case http.StatusInternalServerError:
					if err := ctx.Render("500.gohtml", nil); err != nil {
						b.logFunc(err.Error())
						return
					}
				}
				//ctx.RespData = respData
			}
		}
	}
}
