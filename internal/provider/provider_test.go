package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
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

func testAccPreCheck(t *testing.T) {
	require.NotEmpty(t, os.Getenv("UPTIME_TOKEN"), "UPTIME_TOKEN must be set for acceptance tests")
}
