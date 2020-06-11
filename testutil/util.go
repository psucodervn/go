package testutil

import (
	"os"
	"testing"
)

func SkipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func OnlyIntegrationTest(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTING") == "" {
		t.Skip("Skipping testing when not in integration testing mode")
	}
}
