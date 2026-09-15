package sessions

import (
	"net/http"
	"time"
)

type CookieConfig struct {
	Secure bool
}

func (c CookieConfig) name() string {
	if c.Secure {
		return "__Host-fick_session"
	}

	return "fick_session"
}

func (c CookieConfig) Set(
	w http.ResponseWriter,
	token string,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     c.name(),
			Value:    token,
			Path:     "/",
			Secure:   c.Secure,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	w.Header().Set(
		"Cache-Control",
		"no-store",
	)
}

func (c CookieConfig) Clear(
	w http.ResponseWriter,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     c.name(),
			Value:    "",
			Path:     "/",
			Secure:   c.Secure,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
			Expires:  time.Unix(1, 0),
		},
	)

	w.Header().Set(
		"Cache-Control",
		"no-store",
	)
}

func (c CookieConfig) Name() string {
	if c.Secure {
		return "__Host-fick_session"
	}

	return "fick_session"
}
