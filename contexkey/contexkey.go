package contexkey

import "context"

var (
	emailKey  = contextKey("email")
	tokenKey  = contextKey("token")
	apiurlKey = contextKey("apiurl")
)

type contextKey string

func (c contextKey) String() string {
	return "stock-miner-" + string(c)
}

func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}

func EmailFromContext(ctx context.Context) (string, bool) {
	str, ok := ctx.Value(emailKey).(string)

	return str, ok
}

func WithToken(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, tokenKey, s)
}

func TokenFromContext(ctx context.Context) (string, bool) {
	str, ok := ctx.Value(tokenKey).(string)

	return str, ok
}

func WithAPIURL(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, apiurlKey, s)
}

func APIURLFromContext(ctx context.Context) (string, bool) {
	str, ok := ctx.Value(apiurlKey).(string)

	return str, ok
}
