package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	xymcp "github.com/gumieri/xyops-integrations/internal/mcp"
	"github.com/gumieri/xyops-integrations/sdk"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// version, commit, and buildTime are set at build time by GoReleaser via
// -X main.version, -X main.commit, -X main.buildTime. When building directly
// with `go build`, they fall back to the defaults below.
// commit and buildTime are not consumed by the server itself, but must
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
	baseURL := os.Getenv("XYOPS_BASE_URL")
	if baseURL == "" {
		log.Fatal("XYOPS_BASE_URL environment variable is required")
	}

	apiKey := os.Getenv("XYOPS_API_KEY")
	if apiKey == "" {
		log.Fatal("XYOPS_API_KEY environment variable is required")
	}
	if len(apiKey) < 16 {
		log.Fatal("XYOPS_API_KEY must be at least 16 characters long")
	}

	xyopsClient, err := sdk.New(baseURL, apiKey)
	if err != nil {
		log.Fatalf("Failed to create xyops client: %v", err)
	}

	xymcp.Version = version
	server := xymcp.NewServer(xyopsClient)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, &mcpsdk.StdioTransport{}); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		fmt.Fprintln(os.Stderr, "Server failed:", err)
		os.Exit(1)
	}
}
