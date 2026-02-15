package service

import "testing"

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	c1, err := generateCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c2, err := generateCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c1 == "" || c2 == "" {
		t.Fatal("code should not be empty")
	}

	if len(c1) < 40 {
		t.Fatalf("code too short: %d", len(c1))
	}

	if c1 == c2 {
		t.Fatal("codes should be unique")
	}
}
