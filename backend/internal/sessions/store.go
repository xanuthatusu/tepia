package sessions

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func Init() {
	key := []byte(os.Getenv("SESSION_KEY"))
	if len(key) == 0 {
		panic("SESSION_KEY environment variable is not set")
	}

	Store = sessions.NewCookieStore(key)
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
