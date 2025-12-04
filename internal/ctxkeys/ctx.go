package ctxkeys

import (
	"context"

	"github.com/axadrn/axeladrian/internal/config"
)

type contextKey string

const (
	ConfigKey  contextKey = "config"
	URLPathKey contextKey = "url_path"
)

func Config(ctx context.Context) *config.Config {
	cfg, _ := ctx.Value(ConfigKey).(*config.Config)
	return cfg
}

func WithConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, ConfigKey, cfg)
}

func URLPath(ctx context.Context) string {
	path, _ := ctx.Value(URLPathKey).(string)
	return path
}

func WithURLPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, URLPathKey, path)
}
