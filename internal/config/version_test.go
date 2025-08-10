package config

import "testing"

func TestUserAgent(t *testing.T) {
	oldV, oldC, oldD := Version, Commit, BuildDate
	t.Cleanup(func() { Version, Commit, BuildDate = oldV, oldC, oldD })

	cases := []struct {
		name     string
		version  string
		commit   string
		build    string
		expected string
	}{
		{"defaults", "dev", "none", "unknown", "idpproxy/dev (commit=none)"},
		{"overridden", "1.0.0", "abc1234", "2025-08-10", "idpproxy/1.0.0 (commit=abc1234)"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			Version, Commit, BuildDate = test.version, test.commit, test.build
			if got := UserAgent(); got != test.expected {
				t.Fatalf("UserAgent() = %q, want %q", got, test.expected)
			}
		})
	}
}
