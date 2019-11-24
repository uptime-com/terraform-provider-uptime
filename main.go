package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/uptime-com/terraform-provider-uptime/uptime"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func () terraform.ResourceProvider {
			return uptime.Provider()
		},
	})
}
