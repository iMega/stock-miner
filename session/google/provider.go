package google

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/imega/stock-miner/broker"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	googleApi "google.golang.org/api/oauth2/v2"
)

func GoogleSignInHandlers(
	clientID,
	clientSecret,
	callbackUrl string,
	state gologin.CookieConfig,
	issueSession http.Handler,
) (http.Handler, http.Handler) {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackUrl,
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	loginHandler := google.StateHandler(
		state,
		google.LoginHandler(conf, nil),
	)
	callbackHandler := google.StateHandler(
		state,
		google.CallbackHandler(conf, issueSession, nil),
	)

	return loginHandler, callbackHandler
}

func UserFromContext(ctx context.Context) (broker.User, error) {
	googleUser, err := google.UserFromContext(ctx)
	if err != nil {
		return broker.User{}, err
	}

	return broker.User{
		Email:  googleUser.Email,
		ID:     googleUser.Id,
		Name:   googleUser.Name,
		Avatar: googleUser.Picture,
	}, nil
}

func WithFakeUser(ctx context.Context) context.Context {
	return google.WithUser(ctx, &googleApi.Userinfo{
		Email: "irvis@imega.ru",
	})
}
