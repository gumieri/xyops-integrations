package main

import (
	"context"
	"flag"
	"log"

	xyopsprovider "github.com/gumieri/xyops-integrations/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate terraform fmt -recursive examples/

// version is set at build time by GoReleaser via -X main.version.
// When building directly with `go build`, it falls back to "dev".
var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/gumieri/xyops",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return xyopsprovider.New(version)
	}, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
