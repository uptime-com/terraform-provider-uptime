package uptime

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
)

// Config defines configuration options for the Uptime.com client
type Config struct {
	// Uptime.com API token
	Token string
	RateMilliseconds int
}

const badCredentials = `

No credentials found for Uptime.com provider.
Please provide an API token in the provider block.
`

func (c *Config) Client() (*uptime.Client, error) {
	if c.Token == "" {
		return nil, fmt.Errorf(badCredentials)
	}

	var httpClient *http.Client
	httpClient = http.DefaultClient
	httpClient.Transport = logging.NewTransport("Uptime.com", http.DefaultTransport)

	config := &uptime.Config{
		HTTPClient: httpClient,
		Token: c.Token,
		UserAgent: fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, terraform.VersionString()),
		RateMilliseconds: c.RateMilliseconds,
	}

	client, err := uptime.NewClient(config)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Uptime.com client configured")

	return client, nil
}
