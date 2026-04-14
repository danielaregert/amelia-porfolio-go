package adminsession

import (
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	CookieName  = "admin_session"
	TokenMaxAge = 7 * 24 * 60 * 60
	TokenMaxDur = 7 * 24 * time.Hour
)

func Authenticate(app *pocketbase.PocketBase, email, password string) (*core.Record, string, error) {
	rec, err := app.FindAuthRecordByEmail("_superusers", email)
	if err != nil || rec == nil {
		return nil, "", http.ErrNoCookie
	}
	if !rec.ValidatePassword(password) {
		return nil, "", http.ErrNoCookie
	}
	token, err := rec.NewStaticAuthToken(TokenMaxDur)
	if err != nil {
		return nil, "", err
	}
	return rec, token, nil
}

func SetCookie(w http.ResponseWriter, r *http.Request, token string) {
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   TokenMaxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(TokenMaxDur),
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func CurrentUser(app *pocketbase.PocketBase, r *http.Request) *core.Record {
	c, err := r.Cookie(CookieName)
	if err != nil || c == nil || c.Value == "" {
		return nil
	}
	rec, err := app.FindAuthRecordByToken(c.Value, core.TokenTypeAuth)
	if err != nil || rec == nil {
		return nil
	}
	if rec.Collection().Name != "_superusers" {
		return nil
	}
	return rec
}
