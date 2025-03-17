package rapidgo

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	handlers   []func(c *Context)
	handlerIdx int
	values     map[string]any
	params     map[string]string
}

func NewContext(responseWriter http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Request:    request,
		Response:   responseWriter,
		handlers:   make([]func(c *Context), 0),
		handlerIdx: -1,
		values:     make(map[string]any),
		params:     make(map[string]string),
	}
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}

func (c *Context) JSON(status int, v any) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	json.NewEncoder(c.Response).Encode(v)
}

func (c *Context) Send(v string) {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(200)
	c.Response.Write([]byte(v))
}

func (c *Context) SendStatus(status int) {
	c.Response.WriteHeader(status)
}

func (c *Context) Bind(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

func (c *Context) BindJSON(v any) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Next() {
	c.handlerIdx++
	if c.handlerIdx < len(c.handlers) {
		c.handlers[c.handlerIdx](c)
	}
}

func (c *Context) Abort() {
	c.handlerIdx = len(c.handlers)
}

func (c *Context) AbortWithStatus(status int) {
	c.Response.WriteHeader(status)
	c.handlerIdx = len(c.handlers)
}

func (c *Context) AbortWithStatusJSON(status int, v any) {
	c.JSON(status, v)
	c.handlerIdx = len(c.handlers)
}

func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// Set any key value pair in the context
func (c *Context) Set(key string, value any) {
	c.values[key] = value
}

// Get any key value pair from the context
func (c *Context) Get(key string) any {
	return c.values[key]
}

func (c *Context) Param(key string) string {
	return c.params[key]
}

func (c *Context) Params() map[string]string {
	return c.params
}

func (c *Context) GetCookie(key string) string {
	cookie, err := c.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (c *Context) SetCookie(key, value string, age time.Duration, secure bool) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   int(age.Seconds()),
		Secure:   secure,
		HttpOnly: true,
	}
	c.Response.Header().Set("Set-Cookie", cookie.String())
}

func (c *Context) DeleteCookie(key string) {
	cookie := &http.Cookie{
		Name:   key,
		Value:  "",
		MaxAge: -1,
	}
	c.Response.Header().Set("Set-Cookie", cookie.String())
}
