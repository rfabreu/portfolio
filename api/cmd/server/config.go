package main

import (
	"fmt"
	"strconv"
	"time"
)

const (
	defaultChatRPM         = 10
	defaultGameRPM         = 60
	defaultChatClientDaily = 50
	defaultChatDailyCap    = 500
)

type config struct {
	allowedOrigins   string
	chatRPM          int
	gameRPM          int
	chatClientDaily  int
	chatDailyCap     int
	trustFlyClientIP bool
	// now is the clock for the limiters and the daily cap; nil means time.Now. It exists for tests.
	now func() time.Time
}

// loadConfig reads the rate-limit settings. An invalid value never stops the
// server: it falls back to its default and adds a warning for the caller to
// log. A bad deploy setting should not take the whole site down, and every
// default is a restrictive value, so falling back is the safe direction.
//
// RATE_LIMIT_RPM is deprecated. Before this limiter it was one bucket shared by
// every route and, behind Fly's proxy, probably every visitor, so it says little
// about a sensible per-client value. When set, it is the default for the game
// route; the chat default becomes the smaller of it and the built-in default, so
// it can tighten chat but never loosen it. CHAT_RATE_LIMIT_RPM and
// GAME_RATE_LIMIT_RPM override it.
func loadConfig(getenv func(string) string) (config, []string) {
	var warnings []string
	positive := func(key string, def int) int {
		raw := getenv(key)
		if raw == "" {
			return def
		}
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			warnings = append(warnings, fmt.Sprintf("invalid %s=%q (want a positive integer); using default %d", key, raw, def))
			return def
		}
		return n
	}

	chatDefault, gameDefault := defaultChatRPM, defaultGameRPM
	if raw := getenv("RATE_LIMIT_RPM"); raw != "" {
		if n, err := strconv.Atoi(raw); err != nil || n <= 0 {
			warnings = append(warnings, fmt.Sprintf("invalid RATE_LIMIT_RPM=%q (want a positive integer); ignoring it", raw))
		} else {
			gameDefault = n
			chatDefault = min(n, defaultChatRPM)
			warnings = append(warnings, fmt.Sprintf("RATE_LIMIT_RPM is deprecated; it now sets the game default (%d) and caps the chat default (%d). Use CHAT_RATE_LIMIT_RPM and GAME_RATE_LIMIT_RPM", gameDefault, chatDefault))
		}
	}

	cfg := config{
		chatRPM:         positive("CHAT_RATE_LIMIT_RPM", chatDefault),
		gameRPM:         positive("GAME_RATE_LIMIT_RPM", gameDefault),
		chatClientDaily: positive("CHAT_CLIENT_DAILY_LIMIT", defaultChatClientDaily),
		chatDailyCap:    positive("CHAT_DAILY_CAP", defaultChatDailyCap),
	}
	if raw := getenv("TRUST_FLY_CLIENT_IP"); raw != "" {
		trust, err := strconv.ParseBool(raw)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("invalid TRUST_FLY_CLIENT_IP=%q (want 1 or 0); not trusting Fly-Client-IP", raw))
		}
		cfg.trustFlyClientIP = trust
	}
	return cfg, warnings
}
