package provider

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/uptime-com/terraform-provider-uptime/internal/uptimeapi"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	os.Exit(m.Run())
}

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"uptime": providerserver.NewProtocol6WithError(VersionFactory("test")()),
	}
}

func testAccAPIClient(t *testing.T) uptimeapi.ClientWithResponsesInterface {
	token := os.Getenv("UPTIME_TOKEN")
	require.NotEmpty(t, token, "UPTIME_TOKEN must be set for acceptance tests")

	c, err := uptimeapi.NewClientWithResponses("https://uptime.com", uptimeapi.WithToken(token))
	require.NoError(t, err)

	return c
}

func testAccSetupContactGroup(t *testing.T, api uptimeapi.ClientWithResponsesInterface) {
	ctx, cancel := context.WithCancel(context.Background())

	res, err := api.PostContactgrouplistWithResponse(ctx, uptimeapi.PostContactgrouplistJSONRequestBody{
		Name:      "void",
		EmailList: ptr([]openapi_types.Email{"noreply@uptime.com"}),
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode(), prettyResponse(res.Status(), bytes.NewReader(res.Body)))

	t.Cleanup(func() {
		defer cancel()
		_, err = api.DeleteContactGroupDetailWithResponse(ctx, strconv.Itoa(*res.JSON200.Pk))
		require.NoError(t, err)

	})
}
