package typesafe

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// LevelOff silences SDK logging.
const LevelOff = slog.LevelError + 4

// LogLevels are the names accepted by TYPESAFE_LOG_LEVEL, from most to least verbose.
var LogLevels = []string{"debug", "info", "warn", "error", "off"}

// ParseLogLevel resolves a TYPESAFE_LOG_LEVEL name to an [slog.Level].
func ParseLogLevel(name string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "off":
		return LevelOff, nil
	}
	return 0, newError("invalid log level %q; expected one of: %s", name, strings.Join(LogLevels, ", "))
}

// defaultLogger writes to standard error at the TYPESAFE_LOG_LEVEL level, or at
// warn — above everything this SDK emits — when the variable is unset.
func defaultLogger() (*slog.Logger, error) {
	level := slog.LevelWarn
	if name := env(LogLevelEnv); name != "" {
		parsed, err := ParseLogLevel(name)
		if err != nil {
			return nil, err
		}
		level = parsed
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler).With(slog.String("sdk", sdkName)), nil
}

// keyHeaders keep their scheme and a suffix of the secret so a key can still be
// told apart in a log; opaqueHeaders lose their value entirely.
var (
	keyHeaders    = map[string]bool{"Authorization": true, "Proxy-Authorization": true, "X-Api-Key": true, "Api-Key": true}
	opaqueHeaders = map[string]bool{"Cookie": true, "Set-Cookie": true}
)

func redactHeader(header http.Header) map[string]string {
	redacted := make(map[string]string, len(header))
	for name, values := range header {
		value := strings.Join(values, ", ")
		switch {
		case keyHeaders[name]:
			redacted[name] = redactKey(value)
		case opaqueHeaders[name]:
			redacted[name] = "***"
		default:
			redacted[name] = value
		}
	}
	return redacted
}

func redactKey(value string) string {
	scheme, secret := "", value
	if before, after, found := strings.Cut(value, " "); found {
		scheme, secret = before+" ", strings.TrimSpace(after)
	}
	tail := ""
	if len(secret) > 8 {
		tail = secret[len(secret)-4:]
	}
	return scheme + "***" + tail
}
