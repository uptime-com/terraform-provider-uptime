package uptimeapi

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

type httpClientMock struct {
	mock.Mock
}

func (m *httpClientMock) Do(req *http.Request) (*http.Response, error) {
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

	d := &httpClientMock{}
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

func TestWithToken(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := &httpClientMock{}
	defer d.AssertExpectations(t)

	c, err := NewClient("https://example.com", WithHTTPClient(d), WithToken("foo-bar-baz"))
	if err != nil {
		t.Fatal(err)
	}

	d.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		r := args.Get(0).(*http.Request)
		require.Equal(t, "Token foo-bar-baz", r.Header.Get("Authorization"))
	}).Return(&http.Response{}, nil).Once()

	_, err = c.GetAuthMe(ctx)
	require.NoError(t, err)
}
