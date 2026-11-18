package main

import (
	"strings"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadConfig_Defaults(t *testing.T) {
	cfg, warns := loadConfig(env(nil))
	if cfg.chatRPM != 10 || cfg.gameRPM != 60 || cfg.chatClientDaily != 50 || cfg.chatDailyCap != 500 || cfg.trustFlyClientIP {
		t.Errorf("defaults = %+v", cfg)
	}
	if len(warns) != 0 {
		t.Errorf("unexpected warnings: %v", warns)
	}
}

func TestLoadConfig_ReadsValues(t *testing.T) {
	cfg, warns := loadConfig(env(map[string]string{
		"CHAT_RATE_LIMIT_RPM": "5", "GAME_RATE_LIMIT_RPM": "100", "CHAT_CLIENT_DAILY_LIMIT": "7", "CHAT_DAILY_CAP": "42", "TRUST_FLY_CLIENT_IP": "1",
	}))
	if cfg.chatRPM != 5 || cfg.gameRPM != 100 || cfg.chatClientDaily != 7 || cfg.chatDailyCap != 42 || !cfg.trustFlyClientIP {
		t.Errorf("cfg = %+v", cfg)
	}
	if len(warns) != 0 {
		t.Errorf("unexpected warnings: %v", warns)
	}
}

func TestLoadConfig_InvalidValuesFallBackToDefaultWithWarning(t *testing.T) {
	bad := []string{"abc", "0", "-5", "1.5", " 7", "99999999999999999999"}
	for _, v := range bad {
		for _, key := range []string{"CHAT_RATE_LIMIT_RPM", "GAME_RATE_LIMIT_RPM", "CHAT_CLIENT_DAILY_LIMIT", "CHAT_DAILY_CAP", "RATE_LIMIT_RPM"} {
			cfg, warns := loadConfig(env(map[string]string{key: v}))
			if cfg.chatRPM != 10 || cfg.gameRPM != 60 || cfg.chatClientDaily != 50 || cfg.chatDailyCap != 500 {
				t.Errorf("%s=%q: cfg = %+v, want defaults", key, v, cfg)
			}
			if len(warns) == 0 || !strings.Contains(warns[0], key) {
				t.Errorf("%s=%q: warnings = %v, want one naming the variable", key, v, warns)
			}
		}
	}
}

func TestLoadConfig_InvalidTrustFlagIsOffWithWarning(t *testing.T) {
	cfg, warns := loadConfig(env(map[string]string{"TRUST_FLY_CLIENT_IP": "yes please"}))
	if cfg.trustFlyClientIP {
		t.Error("an unparseable TRUST_FLY_CLIENT_IP must leave trust off")
	}
	if len(warns) != 1 || !strings.Contains(warns[0], "TRUST_FLY_CLIENT_IP") {
		t.Errorf("warnings = %v", warns)
	}
}

func TestLoadConfig_LegacyRateLimitRPMSetsGameAndCanOnlyTightenChat(t *testing.T) {
	cases := []struct {
		legacy     string
		chat, game int
	}{
		{"30", 10, 30}, // a higher legacy value must not loosen chat
		{"600", 10, 600},
		{"5", 5, 5}, // a lower one tightens both
	}
	for _, c := range cases {
		cfg, warns := loadConfig(env(map[string]string{"RATE_LIMIT_RPM": c.legacy}))
		if cfg.chatRPM != c.chat || cfg.gameRPM != c.game {
			t.Errorf("RATE_LIMIT_RPM=%s: cfg = %+v, want chat %d game %d", c.legacy, cfg, c.chat, c.game)
		}
		if len(warns) != 1 || !strings.Contains(warns[0], "deprecated") {
			t.Errorf("RATE_LIMIT_RPM=%s: warnings = %v, want a deprecation notice", c.legacy, warns)
		}
	}
}

func TestLoadConfig_PerRouteWinsOverLegacy(t *testing.T) {
	cfg, _ := loadConfig(env(map[string]string{"RATE_LIMIT_RPM": "30", "CHAT_RATE_LIMIT_RPM": "40", "GAME_RATE_LIMIT_RPM": "2"}))
	if cfg.chatRPM != 40 || cfg.gameRPM != 2 {
		t.Errorf("cfg = %+v, want chat 40 and game 2", cfg)
	}
}

func TestLoadConfig_InvalidPerRouteFallsBackToWhatItWouldOtherwiseUse(t *testing.T) {
	cfg, warns := loadConfig(env(map[string]string{"RATE_LIMIT_RPM": "30", "CHAT_RATE_LIMIT_RPM": "oops"}))
	if cfg.chatRPM != 10 || cfg.gameRPM != 30 {
		t.Errorf("cfg = %+v, want chat 10 and game 30", cfg)
	}
	found := false
	for _, w := range warns {
		found = found || strings.Contains(w, "CHAT_RATE_LIMIT_RPM")
	}
	if !found {
		t.Errorf("warnings = %v, want one naming CHAT_RATE_LIMIT_RPM", warns)
	}
}
