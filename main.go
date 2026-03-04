package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/winebarrel/terraform-provider-statuspage/internal/provider"
)

// Provider documentation generation.
//go:generate go tool tfplugindocs generate --provider-name statuspage

var version string = "dev"

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/winebarrel/statuspage",
		Debug:   *debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
