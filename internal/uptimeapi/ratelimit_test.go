package uptimeapi

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

type httpRequestDoerMock struct {
	mock.Mock
}

func (m *httpRequestDoerMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	return args.Get(0).(*http.Response), nil
}

func TestWithRateLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := &httpRequestDoerMock{}
	defer d.AssertExpectations(t)

	c, err := NewClient("https://example.com", WithHTTPClient(d), WithRateLimitEvery(100*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	started := time.Now()

	d.On("Do", mock.Anything).Return(&http.Response{}, nil).Times(4)
	for i := 0; i < 4; i++ {
		_, err := c.GetAuthMe(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	require.GreaterOrEqual(t, time.Since(started), 300*time.Millisecond)
}
