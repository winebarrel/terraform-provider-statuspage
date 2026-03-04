package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
	"github.com/winebarrel/terraform-provider-statuspage/internal/provider"
)

// Provider documentation generation.
//go:generate go tool tfplugindocs generate --provider-name statuspage

var version string = "dev"

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	apiclientOpts := []apiclient.ClientOption{}
	if os.Getenv("TF_LOG") == "debug" {
		apiclientOpts = append(apiclientOpts, apiclient.WithDebug())
	}

	err := providerserver.Serve(context.Background(), provider.New(version, apiclientOpts...), providerserver.ServeOpts{
		Address: "registry.terraform.io/winebarrel/statuspage",
		Debug:   *debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
