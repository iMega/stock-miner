package contexkey

import "context"

var (
	emailKey = contextKey("email")
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
