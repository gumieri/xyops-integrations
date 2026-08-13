package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	xymcp "github.com/gumieri/xyops-integrations/internal/mcp"
	"github.com/gumieri/xyops-integrations/sdk"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// version is set at build time by GoReleaser via -X main.version.
// When building directly with `go build`, it falls back to "dev".
var version = "dev"

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

	if strings.HasPrefix(baseURL, "http://") {
		log.Printf("WARNING: XYOPS_BASE_URL uses plain HTTP; the API key will be sent in cleartext")
	}
	log.Printf("xyops-mcp %s starting", version)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, &mcpsdk.StdioTransport{}); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Fatalf("Server failed: %v", err)
	}
}
