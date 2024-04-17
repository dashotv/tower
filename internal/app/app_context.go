package app

import "context"

func ContextSet(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func ContextApp(ctx context.Context) *Application {
	return ctx.Value("app").(*Application)
}
