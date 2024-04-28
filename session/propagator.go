package session

import "net/http"

// Propagator 负责将session id 从http request中提取，注入到http response中
type Propagator interface {
	// Inject 将session id注入到http response中
	Inject(id string, writer http.ResponseWriter) error
	// Extract 从http request中提取session id
	Extract(req *http.Request) (string, error)
	// Clean 清除http response中的session id
	Clean(writer http.ResponseWriter) error
}

type CookiePropagatorOption func(propagator *CookiePropagator)

// CookiePropagator 基于cookie的session id传递器
type CookiePropagator struct {
	cookieName   string
	cookieOption func(cookie *http.Cookie)
}

func NewCookiePropagator() Propagator {
	return &CookiePropagator{
		cookieName: "session_id",
		cookieOption: func(cookie *http.Cookie) {

		},
	}
}

func WithCookieName(name string) CookiePropagatorOption {
	return func(propagator *CookiePropagator) {
		propagator.cookieName = name
	}
}

func WithCookieOption(option func(cookie *http.Cookie)) CookiePropagatorOption {
	return func(propagator *CookiePropagator) {
		propagator.cookieOption = option
	}
}

func (c *CookiePropagator) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  c.cookieName,
		Value: id,
		Path:  "/",
	}
	c.cookieOption(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (c *CookiePropagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(c.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (c *CookiePropagator) Clean(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   c.cookieName,
		MaxAge: -1,
	}
	http.SetCookie(writer, cookie)
	return nil
}
