package anonymous

import (
	"context"
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
)

var apiContextKey = config.Repository + "/" + fmt.Sprintf("%T", anonymous{})

func Context(ctx context.Context, an Anonymous) context.Context {
	ctx = context.WithValue(ctx, apiContextKey, an)
	return ctx
}

func FromContext(ctx context.Context) Anonymous {
	return ctx.Value(apiContextKey).(Anonymous)
}
