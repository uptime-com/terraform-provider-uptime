package main

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.13.0
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,client,spec -package uptimeapi -o internal/uptimeapi/uptimeapi.gen.go https://uptime.com/api/v1/openapi/

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/uptime-com/terraform-provider-uptime/internal/provider"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/uptime-com/uptime",
		Debug:           false,
		ProtocolVersion: 5,
	}

	flag.BoolVar(&opts.Debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	log.Printf("terraform-provider-uptime %s, commit %s, built at %s", version, commit, date)

	err := providerserver.Serve(context.Background(), provider.VersionFactory(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
