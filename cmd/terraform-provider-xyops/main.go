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

// version, commit, and buildTime are set at build time by GoReleaser via
// -X main.version, -X main.commit, -X main.buildTime. When building directly
// with `go build`, they fall back to the defaults below.
// commit and buildTime are not consumed by the provider itself, but must
// remain declared so GoReleaser's -X ldflags resolve; the init references
// keep the `unused` linter happy.
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func init() {
	_ = commit
	_ = buildTime
}

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
