package lr

import (
	"bytes"
	"context"
	"html/template"
)

// TemplateEngine 模版引擎接口
type TemplateEngine interface {
	// Render 渲染页面
	// tplName 模版的名字
	// data 需要渲染的页面数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

// GoTemplateEngine 基于go的基础包实现的模版引擎
type GoTemplateEngine struct {
	T *template.Template
}

func (g GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	if err := g.T.ExecuteTemplate(bs, tplName, data); err != nil {
		return nil, err
	}
	return bs.Bytes(), nil
}
