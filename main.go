package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/uptime-com/terraform-provider-uptime/uptime"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: uptime.Provider,
	})
}
