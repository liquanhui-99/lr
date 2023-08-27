package lr

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

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
