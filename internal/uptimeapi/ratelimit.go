package uptimeapi

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func WithRateLimit(r rate.Limit) ClientOption {
	return withRateLimiter(rate.NewLimiter(r, 1))
}

func WithRateLimitEvery(duration time.Duration) ClientOption {
	return withRateLimiter(rate.NewLimiter(rate.Every(duration), 1))
}

func withRateLimiter(l *rate.Limiter) ClientOption {
	return func(c *Client) error {
		c.Client = &rateLimitedHttpRequestDoer{
			c: c.Client,
			l: l,
		}
		return nil
	}
}

type rateLimitedHttpRequestDoer struct {
	c HttpRequestDoer
	l *rate.Limiter
}

func (c *rateLimitedHttpRequestDoer) Do(req *http.Request) (*http.Response, error) {
	err := c.l.Wait(req.Context())
	if err != nil {
		return nil, err
	}
	return c.c.Do(req)
}
