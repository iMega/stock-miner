package session

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/containerd/containerd/log"
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/sessions"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/session/google"
)

const (
	sessionName  = "stock-miner"
	cookieMaxAge = 60
)

type SessionStore struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	db           *sessions.CookieStore
	userDB       domain.UserStorage
	isDevMode    bool
	RootEmail    string
}

func New(opts ...Option) *SessionStore {
	s := &SessionStore{}

	for _, opt := range opts {
		opt(s)
	}

	if s.isDevMode {
		ctx := contexkey.WithEmail(context.Background(), s.RootEmail)
		s.userDB.CreateUser(ctx, domain.User{
			ID:    "1",
			Email: s.RootEmail,
			Role:  "root",
		})
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

func WithUserStorage(s domain.UserStorage) Option {
	return func(p *SessionStore) {
		p.userDB = s
	}
}

func WithDevMode(s string) Option {
	return func(p *SessionStore) {
		if s == "true" {
			p.isDevMode = true
		}
	}
}

func WithRootEmail(s string) Option {
	return func(p *SessionStore) {
		p.RootEmail = s
	}
}

func (s *SessionStore) AppendHandlers(mux *http.ServeMux) {
	login, callback := google.GoogleSignInHandlers(
		s.ClientID,
		s.ClientSecret,
		s.CallbackURL,
		gologin.CookieConfig{
			Name:     "stock-miner-tmp",
			Path:     "/",
			MaxAge:   cookieMaxAge,
			HTTPOnly: true,
			Secure:   !s.isDevMode, // HTTP only dev
		},
		s.issueSession(),
	)

	if s.isDevMode {
		login = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/google/callback", http.StatusFound)
		})

		callback = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			r := req.WithContext(google.WithFakeUser(req.Context()))
			s.issueSession().ServeHTTP(w, r)
		})
	}

	mux.Handle("/google/login", login)
	mux.Handle("/google/callback", callback)
	mux.Handle("/logout", s.logoutHandler())
}

func (s *SessionStore) DefenceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.db.Get(r, sessionName)
		if err != nil {
			r.URL.Path = "/signin.htm"
			next.ServeHTTP(w, r)

			return
		}

		if strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = "/index.htm"
		}

		email, _ := session.Values["email"].(string)
		ctx := contexkey.WithEmail(r.Context(), email)
		next.ServeHTTP(w, r.WithContext(ctx))
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
		user, err := google.UserFromContext(ctx)
		if err != nil {
			log.GetLogger(ctx).Errorf("failed to extract user from context, %s", err)
			http.Error(w, ":(", http.StatusInternalServerError)
			return
		}

		session := s.db.New(sessionName)

		session.Values["id"] = user.ID
		session.Values["name"] = user.Name
		session.Values["email"] = user.Email

		if err := session.Save(w); err != nil {
			log.GetLogger(ctx).Errorf("failed to save session, %s", err)
			http.Error(w, ":(", http.StatusInternalServerError)

			return
		}

		ctxNew := contexkey.WithEmail(ctx, user.Email)
		if _, err = s.userDB.GetUser(ctxNew); err != nil {
			user.Role = "user"
			if user.Email == s.RootEmail {
				user.Role = "root"
			}
			if err := s.userDB.CreateUser(ctxNew, user); err != nil {
				log.GetLogger(ctx).Errorf("failed to create user, %s", err)
				http.Error(w, ":(", http.StatusInternalServerError)

				return
			}
		}

		http.Redirect(w, req.WithContext(ctxNew), "/", http.StatusFound)
	})
}
