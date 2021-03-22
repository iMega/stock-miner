package session_store

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/dghubble/sessions"
	provider_google "github.com/imega/stock-miner/session_store/google"
)

const sessionName = "stock-miner"

type SessionStore struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	db           *sessions.CookieStore
}

func New(opts ...Option) *SessionStore {
	s := &SessionStore{}

	for _, opt := range opts {
		opt(s)
	}

	s.db = sessions.NewCookieStore([]byte(s.ClientSecret), nil)

	return s
}

type Option func(p *SessionStore)

func WithClintID(ID string) Option {
	return func(p *SessionStore) {
		p.ClientID = ID
	}
}

func WithClientSecret(s string) Option {
	return func(p *SessionStore) {
		p.ClientSecret = s
	}
}

func WithCallbackURL(s string) Option {
	return func(p *SessionStore) {
		p.CallbackURL = s
	}
}

func (s *SessionStore) AppendHandlers(mux *http.ServeMux) {
	login, callback := provider_google.GoogleSignInHandlers(
		s.ClientID,
		s.ClientSecret,
		s.CallbackURL,
		gologin.CookieConfig{
			Name:     "stock-miner-tmp",
			Path:     "/",
			MaxAge:   60, // 60 seconds
			HTTPOnly: true,
			Secure:   false, // HTTPS only
		},
		s.issueSession(),
	)

	mux.Handle("/google/login", login)
	mux.Handle("/google/callback", callback)
	mux.Handle("/logout", s.logoutHandler())
}

func (s *SessionStore) DefenceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := s.db.Get(r, sessionName)
		if err != nil {
			r.URL.Path = "/signin.htm"

			next.ServeHTTP(w, r)
			// w.Write([]byte(`<html><body><a href="/google/login">Login with Google</a></body></html>`))
			// http.Redirect(w, r, "/google/login", http.StatusFound)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = "/index.htm"
		}

		// w.Write([]byte(`<p>You are logged in %s!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`))
		// r.URL.Path = "/index.htm"
		next.ServeHTTP(w, r)
	})
}

func (s *SessionStore) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s.db.Destroy(w, sessionName)

		http.SetCookie(w, &http.Cookie{
			Name:    sessionName + "-tmp",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})

		http.SetCookie(w, &http.Cookie{
			Name:    sessionName,
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})

		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func (s *SessionStore) issueSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, ":(", http.StatusInternalServerError)
			return
		}

		session := s.db.New(sessionName)

		fmt.Printf("===%#v\n", googleUser)

		session.Values["userid"] = googleUser.Id
		session.Values["username"] = googleUser.Name
		session.Values["useremail"] = googleUser.Email
		session.Values["usertype"] = "google"

		session.Save(w)

		http.Redirect(w, req, "/", http.StatusFound)
	})
}
