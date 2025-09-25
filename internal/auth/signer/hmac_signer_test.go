package signer

import "testing"

func TestHMACSigner_InfoMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  []byte
		kid  string
	}{
		{"basic", []byte("k1"), "kid-1"},
		{"empty-kid", []byte("k2"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewHMACSigner(tt.key, tt.kid)
			if got := s.Alg(); got != "HS256" {
				t.Fatalf("Alg() = %q", got)
			}
			if got := s.KeyID(); got != tt.kid {
				t.Fatalf("KeyID() = %q", got)
			}
		})
	}
}
