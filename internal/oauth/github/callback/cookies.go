package callback

import "net/http"

const stateCookieName = "oauth_state"

func deleteStateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
}

func safeCookieVal(c *http.Cookie) string {
	if c == nil || c.Value == "" {
		return ""
	}
	if len(c.Value) <= 6 {
		return c.Value
	}
	return c.Value[:6] + "..."
}
