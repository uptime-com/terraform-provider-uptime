package uptime

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/uptime-com/uptime-client-go"
)

var testAccProtoV5ProviderFactoryMap = map[string]func() (tfprotov5.ProviderServer, error){
	"uptime": func() (tfprotov5.ProviderServer, error) {
		return schema.NewGRPCProviderServer(Provider()), nil
	},
}

func testAccAPIClient(t *testing.T) *uptime.Client {
	token := os.Getenv("UPTIME_TOKEN")
	require.NotEmpty(t, token, "UPTIME_TOKEN must be set for acceptance tests")

	c := Config{
		Token:            token,
		RateMilliseconds: 500,
	}
	api, err := c.Client()
	require.NoError(t, err, "failed to initialize uptime.com api client")

	return api
}
