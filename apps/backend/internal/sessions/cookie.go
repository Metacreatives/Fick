package sessions

import "net/http"

type CookieConfig struct {
	Secure bool
}

func (c CookieConfig) Set(
	w http.ResponseWriter,
	token string,
) {
	name := "fick_session"

	if c.Secure {
		name = "__Host-fick_session"
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     name,
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
