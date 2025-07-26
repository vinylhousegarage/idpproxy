package cookie

import (
	"net/http"
	"time"
)

func SetIDTokenCookie(w http.ResponseWriter, idToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    idToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})
}
