package integration

import (
	"os"
	"testing"
)

// TestConductorVersion records the live conductor's version in the test
// output and enforces an optional pin. Validation results are only
// meaningful against a known xyops version (XYOINT-15): every live run
// logs it, and XYOPS_EXPECT_VERSION=1.1.7 fails the run on mismatch.
func TestConductorVersion(t *testing.T) {
	skipIfNoClient(t)

	masters, err := xyopsClient.Servers.Masters(t.Context())
	if err != nil {
		t.Fatalf("list masters: %v", err)
	}
	if len(masters) == 0 {
		t.Skip("no conductor masters reported")
	}

	version := ""
	for _, m := range masters {
		if m.Version != "" {
			version = m.Version
			break
		}
	}
	if version == "" {
		t.Skip("conductor version not reported")
	}

	t.Logf("conductor version: %s", version)

	if want := os.Getenv("XYOPS_EXPECT_VERSION"); want != "" && want != version {
		t.Errorf("conductor version %s does not match XYOPS_EXPECT_VERSION=%s", version, want)
	}
}
