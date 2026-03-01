package service

import "testing"

func TestGenerateProxyCode(t *testing.T) {
	t.Parallel()

	c1, err := generateProxyCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c2, err := generateProxyCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c1 == "" || c2 == "" {
		t.Fatal("proxycode should not be empty")
	}

	if len(c1) < 40 {
		t.Fatalf("proxycode too short: %d", len(c1))
	}

	if c1 == c2 {
		t.Fatal("generated codes should be unique")
	}
}
