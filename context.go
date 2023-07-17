package lorm

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	// Request
	Req *http.Request
	// Response
	Resp http.ResponseWriter
	// 缓存query中的数据
	queryValues url.Values
	// 缓存请求路径参数
	pathParams map[string]string
}

// BindJson 解析json数据
func (c *Context) BindJson(val any) error {
	if val == nil {
		return errors.New("输入不能为nil")
	}

	if c.Req.Body == nil {
		return errors.New("请求body不能为nil")
	}

	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

// BindForm 解析form类型的参数信息
func (c *Context) BindForm(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{
			val: "",
			err: err,
		}
	}
	val, ok := c.Req.Form[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("key不存在"),
		}
	}
	return StringValue{
		val: val[0],
		err: nil,
	}
}

// QueryValue 获取query中的信息, c.Req.URL.Query()没有缓存query，每次都去parse,
// 所以该方法需要建立缓存，并且区分是key不存在，还是存在但值为空字符串
func (c *Context) QueryValue(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("key不存在"),
		}
	}

	return StringValue{
		val: vals[0],
		err: nil,
	}
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.pathParams[key]
	if !ok {
		return StringValue{
			val: "",
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

// Int64 转换string为int64
func (s *StringValue) Int64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// Uint64 转换string为Uint64
func (s *StringValue) Uint64() (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseUint(s.val, 10, 64)
}

// Int 转换string为int
func (s *StringValue) Int() (int, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.Atoi(s.val)
}

// Uint 转换string为uint
func (s *StringValue) Uint() (uint, error) {
	if s.err != nil {
		return 0, s.err
	}

	res, err := strconv.ParseUint(s.val, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(res), nil
}
