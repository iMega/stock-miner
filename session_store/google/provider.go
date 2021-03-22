package provider_google

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
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
