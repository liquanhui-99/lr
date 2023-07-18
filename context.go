package lorm

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// JSON类型的数据
	jsonType = iota
	// XML类型的数据
	xmlType
	// Protobuf类型的数据
	protobufType
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

// 处理请求数据的方法
//解析Json、XML、Protobuf、Query、Form和Path参数

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

// BindXml 解析xml数据
func (c *Context) BindXml(val any) error {
	if val == nil {
		return errors.New("输入不能为nil")
	}
	if c.Req.Body == nil {
		return errors.New("请求body不能为nil")
	}

	xl := xml.NewDecoder(c.Req.Body)
	return xl.Decode(val)
}

// BindProtobuf 解析protobuf数据
func (c *Context) BindProtobuf(val proto.Message) error {
	if val == nil {
		return errors.New("输入不能为nil")
	}
	if c.Req.Body == nil {
		return errors.New("请求body不能为nil")
	}
	body, err := io.ReadAll(c.Req.Body)
	if err != nil {
		return err
	}

	buf := proto.NewBuffer(body)
	return buf.DecodeMessage(val)
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
func (s StringValue) Int64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// Uint64 转换string为Uint64
func (s StringValue) Uint64() (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseUint(s.val, 10, 64)
}

// Int 转换string为int
func (s StringValue) Int() (int, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.Atoi(s.val)
}

// Uint 转换string为uint
func (s StringValue) Uint() (uint, error) {
	if s.err != nil {
		return 0, s.err
	}

	res, err := strconv.ParseUint(s.val, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(res), nil
}

// Int32 转换string为int32
func (s StringValue) Int32() (int32, error) {
	if s.err != nil {
		return 0, s.err
	}
	res, err := strconv.ParseInt(s.val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(res), nil
}

// Uint32 转换string为uint32
func (s StringValue) Uint32() (uint32, error) {
	if s.err != nil {
		return 0, s.err
	}
	res, err := strconv.ParseUint(s.val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(res), nil
}

// Float32 转换字符串为float32
func (s StringValue) Float32() (float32, error) {
	if s.err != nil {
		return 0, s.err
	}

	res, err := strconv.ParseFloat(s.val, 32)
	if err != nil {
		return 0, err
	}
	return float32(res), nil
}

// Float64 转换字符串为float64
func (s StringValue) Float64() (float64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseFloat(s.val, 64)
}

// 处理返回相应的方法

// bindResp 处理响应的基础方法
func (c *Context) bindResp(status int, val []byte, tp int) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	c.Resp.WriteHeader(status)
	switch tp {
	case jsonType:
		c.Resp.Header().Set("Content-Type", "application/json")
	case xmlType:
		c.Resp.Header().Set("Content-Type", "application/xml")
	case protobufType:
		c.Resp.Header().Set("Content-Type", "application/octet-stream")
	}
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))

	n, err := c.Resp.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("未返回全部的数据")
	}
	return nil
}

// RespJsonOk 成功的Json响应
func (c *Context) RespJsonOk(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.bindResp(http.StatusOK, jsonData, jsonType)
}

// RespXMLOK 成功的XML响应
func (c *Context) RespXMLOK(val any) error {
	xmlData, err := xml.Marshal(val)
	if err != nil {
		return err
	}
	return c.bindResp(http.StatusOK, xmlData, xmlType)
}

// RespProtobufOK 成功的protobuf响应
func (c *Context) RespProtobufOK(val proto.Message) error {
	protoData, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	return c.bindResp(http.StatusOK, protoData, protobufType)
}

// RespOKWithMessage 字符串格式的成功响应
func (c *Context) RespOKWithMessage(val string) {
	c.Resp.WriteHeader(http.StatusOK)
	_, _ = c.Resp.Write([]byte(val))
}
