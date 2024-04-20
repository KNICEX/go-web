package web

import (
	"encoding/json"
	"errors"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
)

const abortIndex int = math.MaxInt

type Context struct {
	Req          *http.Request
	Resp         http.ResponseWriter
	PathParams   map[string]string
	queryCache   url.Values
	MatchedRoute string
	Values       map[string]any

	index    int
	handlers []HandleFunc

	StatusCode int
	RespData   []byte
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Req:    req,
		Resp:   w,
		index:  -1,
		Values: make(map[string]any),
	}
}

func (c *Context) Get(key string) (any, bool) {
	val, ok := c.Values[key]
	return val, ok
}

func (c *Context) Set(key string, val any) {
	c.Values[key] = val
}

func (c *Context) Status(status int) {
	c.StatusCode = status
}

func (c *Context) JSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json")
	c.StatusCode = status
	c.RespData = data
	return nil
}

func (c *Context) String(status int, val string) error {
	c.Resp.Header().Set("Content-Type", "text/plain")
	c.StatusCode = status
	c.RespData = []byte(val)
	return nil
}

func (c *Context) HTML(status int, val string) error {
	c.Resp.Header().Set("Content-Type", "text/html")
	c.StatusCode = status
	c.RespData = []byte(val)
	return nil
}

func (c *Context) JsonOK(val any) error {
	return c.JSON(http.StatusOK, val)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) GetCookie(name string) (*http.Cookie, bool) {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return nil, false
	}
	return cookie, true
}

func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("nil pointer")
	}
	return json.NewDecoder(c.Req.Body).Decode(val)
}

func (c *Context) Param(key string) string {
	return c.PathParams[key]
}

func (c *Context) FormValue(key string) (string, bool) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", false
	}
	vals, ok := c.Req.Form[key]
	if !ok {
		return "", false
	}
	return vals[0], true
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Req.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}
	return c.Req.MultipartForm, nil
}

func (c *Context) QueryValue(key string) (string, bool) {
	if c.queryCache == nil {
		c.queryCache = c.Req.URL.Query()
	}
	vals, ok := c.queryCache[key]
	if !ok {
		return "", false
	}
	return vals[0], true
}

func (c *Context) PathValue(key string) (string, bool) {
	res, ok := c.PathParams[key]
	return res, ok
}

func (c *Context) Next() {
	c.index++
	for n := len(c.handlers); c.index < n; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}
