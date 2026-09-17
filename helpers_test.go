package typesafe

import "testing"

// clearEnv isolates a test from the developer's own TypeSafe configuration.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{APIKeyEnv, BaseURLEnv, DefaultModelEnv, LogLevelEnv} {
		t.Setenv(name, "")
	}
}
