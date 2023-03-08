package uptimeapi

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func WithToken(token string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Token "+token)
		return nil
	})
}

func WithRateLimit(r rate.Limit) ClientOption {
	lim := rate.NewLimiter(r, 1)
	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		return lim.Wait(ctx)
	})
}

func WithRateLimitEvery(duration time.Duration) ClientOption {
	return WithRateLimit(rate.Every(duration))
}
