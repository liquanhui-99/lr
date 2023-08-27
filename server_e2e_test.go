//go:build e2e

package lr

import (
	"html/template"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"testing"
)

func TestServerE2e(t *testing.T) {
	h := NewHTTPServer("tcp", ":8081")

	h.POST("/user/profile", func(ctx *Context) {
		t.Log("成功")
		ctx.Resp.WriteHeader(http.StatusOK)
	})

	h.GET("/user/login/:id", func(ctx *Context) {
		id := ctx.pathParams["id"]
		t.Log("参数为：", id)
	})

	if err := h.Server(); err != nil {
		panic(err)
	}
}

func TestFileUpload(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/upload.gohtml")
	if err != nil {
		panic(err)
	}

	h := NewHTTPServer("tcp", ":8081", Template(&GoTemplateEngine{
		T: tpl,
	}))

	h.GET("/upload", func(ctx *Context) {
		data := struct {
			Name string
		}{
			Name: "my-file",
		}
		if err := ctx.Render("upload.gohtml", data); err != nil {
			t.Log(err)
			ctx.Status = http.StatusInternalServerError
			ctx.RespData = []byte("内部错误")
			return
		}
		ctx.Status = http.StatusOK
	})

	hd := FileUploader{
		FileField: "my-file",
		DestPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", header.Filename)
		},
	}
	h.POST("/upload", hd.Handler())

	fd := FileDownloader{
		Dir: filepath.Join("testdata", "download"),
	}
	h.GET("/download", fd.Handle())

	if err := h.Server(); err != nil {
		panic(err)
	}
}
