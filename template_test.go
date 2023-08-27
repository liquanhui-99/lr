package lr

import (
	"html/template"
	"testing"
)

func TestTemplate(t *testing.T) {
	ptl, err := template.ParseGlob("testdata/404.gohtml")
	if err != nil {
		t.Log(err)
		return
	}

	engine := &GoTemplateEngine{
		T: ptl,
	}

	h := NewHTTPServer("tcp", ":8081", Template(engine))

	h.GET("/user", func(ctx *Context) {
		if err := ctx.Render("404.gohtml", nil); err != nil {
			t.Log(err)
		}
	})

	if err := h.Server(); err != nil {
		panic(err)
	}
}
