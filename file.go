package lr

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileUploader 文件上传功能实现
type FileUploader struct {
	// 上传文件的字段名
	FileField string
	// 用户传递方法来处理文件存储的路径，不再方法里处理，返回值是文件的存储路径
	DestPathFunc func(*multipart.FileHeader) string
}

func (f *FileUploader) Handler() HandleFunc {
	return func(ctx *Context) {
		// 处理上传的逻辑
		file, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.Status = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		// 处理要保存的目标路径，打开文件
		path := f.DestPathFunc(fileHeader)
		newFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			ctx.Status = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}

		// 开始copy
		_, err = io.CopyBuffer(newFile, file, nil)
		if err != nil {
			ctx.Status = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}

		ctx.Status = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

// FileDownloader 文件下载功能实现
type FileDownloader struct {
	Dir string
}

func (f FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		res := ctx.QueryValue("file")
		if res.err != nil {
			ctx.Status = 400
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		// 全路径
		req := filepath.Clean(res.val)
		dest := filepath.Join(f.Dir, req)
		// 这里做路径校验，防止攻击者发送相对路径导致下载其他的内部文件
		dest, err := filepath.Abs(dest)
		if err != nil {
			ctx.Status = 400
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		if !strings.Contains(dest, f.Dir) {
			ctx.Status = 400
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		// 文件名
		fileName := filepath.Base(dest)

		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fileName)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, dest)
	}
}
