package session

import "net/http"

type Propagator interface {
	Inject(id string, writer http.ResponseWriter) error
	Extract(req *http.Request) (string, error)
	Clean(writer http.ResponseWriter) error
}

type CookiePropagatorOption func(propagator *CookiePropagator)

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
