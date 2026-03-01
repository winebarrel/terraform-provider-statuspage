package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/winebarrel/terraform-provider-statuspage/internal/provider"
)

var version string = "dev"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/winebarrel/statuspage",
	})
	if err != nil {
		log.Fatal(err)
	}
}
