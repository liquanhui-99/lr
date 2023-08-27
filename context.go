package lr

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	// Req 接受的请求信息
	Req *http.Request
	// Resp 返回响应
	Resp http.ResponseWriter
	// 路径参数
	pathParams map[string]string
	// Query参数的缓存(?)
	queryCache url.Values
	// 匹配到的路径
	matchedPath string
	// 响应状态码
	Status int
	// 请求响应的数据
	RespData []byte
	// 模版渲染引擎
	TplEngine TemplateEngine
}

func (c *Context) RespJsonOK(val any) error {
	return c.respJson(val, http.StatusOK)
}

func (c *Context) respJson(val any, status int) error {
	if val == nil {
		return errors.New("返回值为nil")
	}

	bytes, err := json.Marshal(&val)
	if err != nil {
		return err
	}

	c.RespData = bytes

	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	c.Status = status

	return nil
}

// BindJson 绑定json
func (c *Context) BindJson(val any) error {
	if val == nil {
		return errors.New("输入不能为nil")
	}

	if c.Req.Body == nil {
		return errors.New("body不能为nil")
	}

	decode := json.NewDecoder(c.Req.Body)
	return decode.Decode(&val)
}

func (c *Context) BindXML(val any) error {
	if val == nil {
		return errors.New("输入不能为nil")
	}

	decode := xml.NewDecoder(c.Req.Body)
	return decode.Decode(&val)
}

// FormValue 获取表单中指定key的内容，多次调用ParseForm()，只会解析一次，不会每次都解析
func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{
			err: err,
		}
	}

	return StringValue{
		val: c.Req.FormValue(key),
		err: nil,
	}
}

// MatchedPath 获取请求最终匹配到的路径
func (c *Context) MatchedPath() string {
	return c.matchedPath
}

// QueryValue 根据key获取Query中的值
func (c *Context) QueryValue(key string) StringValue {
	if c.queryCache == nil {
		c.queryCache = c.Req.URL.Query()
	}

	val, ok := c.queryCache[key]
	if !ok || len(c.queryCache) == 0 {
		return StringValue{
			err: errors.New("key不存在"),
		}
	}

	return StringValue{
		val: val[0],
		err: nil,
	}
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.pathParams[key]
	if !ok || len(c.pathParams) == 0 {
		return StringValue{
			err: errors.New("key不存在"),
		}
	}

	return StringValue{
		val: val,
		err: nil,
	}
}

type StringValue struct {
	val string
	err error
}

// String 获取字符串类型的返回值
func (s StringValue) String() (string, error) {
	return s.val, s.err
}

// Int 获取int类型的返回值
func (s StringValue) Int() (int, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.Atoi(s.val)
}

// Int64 获取int64类型的返回值
func (s StringValue) Int64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseInt(s.val, 10, 64)
}

// Int32 获取int32类型的返回值
func (s StringValue) Int32() (int32, error) {
	if s.err != nil {
		return 0, s.err
	}

	val, err := strconv.ParseUint(s.val, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}

// Uint64 获取uint64类型的返回值
func (s StringValue) Uint64() (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseUint(s.val, 10, 64)
}

// Uint32 获取uint32类型的返回值
func (s StringValue) Uint32() (uint32, error) {
	if s.err != nil {
		return 0, s.err
	}

	val, err := strconv.ParseUint(s.val, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(val), nil
}

// Float32 获取float32类型的返回值
func (s StringValue) Float32() (float32, error) {
	if s.err != nil {
		return 0, s.err
	}

	val, err := strconv.ParseFloat(s.val, 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

// Float64 获取float64类型的返回值
func (s StringValue) Float64() (float64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseFloat(s.val, 32)
}

// Render 模版渲染的方法
// @param  tplName 模版名称
// @param data需要被模版渲染的数据
// @return []byte 模版渲染后的数据
func (c *Context) Render(tplName string, data any) error {
	var err error
	c.RespData, err = c.TplEngine.Render(c.Req.Context(), tplName, data)
	if err != nil {
		c.Status = http.StatusInternalServerError
		return err
	}
	c.Status = http.StatusOK

	return nil
}
