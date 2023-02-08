package uptimeapi

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	os.Exit(m.Run())
}

func TestClient(t *testing.T) {
	if os.Getenv("UPTIME_TOKEN") == "" {
		t.Skip("UPTIME_TOKEN not set")
	}

	client, err := NewClientWithResponses("https://uptime.com", WithToken(os.Getenv("UPTIME_TOKEN")))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("tags", func(t *testing.T) {
		obj1, err := client.PostServicetaglistWithResponse(ctx, PostServicetaglistJSONRequestBody{
			Tag:      "test",
			ColorHex: "000000",
		})
		require.NoError(t, err)
		assert.Equal(t, 200, obj1.StatusCode())

		obj2, err := client.GetServicetaglistWithResponse(ctx, &GetServicetaglistParams{})
		require.NoError(t, err)
		assert.Equal(t, 200, obj2.StatusCode())
	})
}

// TestCleanup is a cleanup utility. It will only run if UPTIME_CLEANUP environment variable is set to any non-empty
// value.
//
// WARNING: Upon successful execution it will wipe out all data in configured Uptime.com account!
func TestCleanup(t *testing.T) {
	if os.Getenv("UPTIME_TOKEN") == "" {
		t.Skip("UPTIME_TOKEN not set")
	}
	if os.Getenv("UPTIME_CLEANUP") == "" {
		t.Skip("UPTIME_CLEANUP not set")
	}

	api, err := NewClientWithResponses("https://uptime.com", WithToken(os.Getenv("UPTIME_TOKEN")))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("checks", func(t *testing.T) {
		res, err := api.GetServicelistWithResponse(ctx, &GetServicelistParams{PageSize: ptr(10000)})
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, 200, res.StatusCode())
		if *res.JSON200.Count == 0 {
			t.Skip("no checks to delete")
		}
		for _, obj := range *res.JSON200.Results {
			_, err := api.DeleteServiceDetail(ctx, strconv.Itoa(*obj.Pk))
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("integrations", func(t *testing.T) {
		res, err := api.GetIntegrationlistWithResponse(ctx, &GetIntegrationlistParams{PageSize: ptr(10000)})
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, 200, res.StatusCode())
		if *res.JSON200.Count == 0 {
			t.Skip("no integrations to delete")
		}
		for _, obj := range *res.JSON200.Results {
			_, err := api.DeleteIntegrationDetail(ctx, strconv.Itoa(*obj.Pk))
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("tags", func(t *testing.T) {
		res, err := api.GetServicetaglistWithResponse(ctx, &GetServicetaglistParams{PageSize: ptr(10000)})
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, 200, res.StatusCode())
		if *res.JSON200.Count == 0 {
			t.Skip("no tags to delete")
		}
		for _, obj := range *res.JSON200.Results {
			_, err := api.DeleteServiceTagDetail(ctx, strconv.Itoa(*obj.Pk))
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func ptr[T any](v T) *T {
	return &v
}
