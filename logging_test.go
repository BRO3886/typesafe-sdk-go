package typesafe

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		want    slog.Level
		wantErr bool
	}{
		{name: "debug", want: slog.LevelDebug},
		{name: "INFO", want: slog.LevelInfo},
		{name: " warn ", want: slog.LevelWarn},
		{name: "error", want: slog.LevelError},
		{name: "off", want: LevelOff},
		{name: "chatty", wantErr: true},
		{name: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseLogLevel(test.name)
			if (err != nil) != test.wantErr {
				t.Fatalf("ParseLogLevel(%q) error = %v, wantErr = %v", test.name, err, test.wantErr)
			}
			if err == nil && got != test.want {
				t.Errorf("ParseLogLevel(%q) = %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestRedactHeader(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer sk-live-0123456789abcdef")
	header.Set("X-Api-Key", "short")
	header.Set("Cookie", "session=abc123")
	header.Set("X-Tenant", "acme")

	redacted := redactHeader(header)

	want := map[string]string{
		"Authorization": "Bearer ***cdef",
		"X-Api-Key":     "***",
		"Cookie":        "***",
		"X-Tenant":      "acme",
	}
	for name, expected := range want {
		if redacted[name] != expected {
			t.Errorf("redactHeader()[%q] = %q, want %q", name, redacted[name], expected)
		}
	}
}

func TestDefaultLoggerReadsTheEnvironment(t *testing.T) {
	clearEnv(t)
	t.Setenv(LogLevelEnv, "debug")

	logger, err := defaultLogger()
	if err != nil {
		t.Fatalf("defaultLogger() error = %v", err)
	}
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("defaultLogger() ignored TYPESAFE_LOG_LEVEL=debug")
	}

	t.Setenv(LogLevelEnv, "off")
	logger, err = defaultLogger()
	if err != nil {
		t.Fatalf("defaultLogger() error = %v", err)
	}
	if logger.Enabled(context.Background(), slog.LevelError) {
		t.Error("defaultLogger() still logged with TYPESAFE_LOG_LEVEL=off")
	}
}
