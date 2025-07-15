package config

import (
	"os"
	"testing"
)

func TestGetPort_WithEnvVar(t *testing.T) {
	os.Setenv("PORT", "12345")
	defer os.Unsetenv("PORT")

	got := GetPort()
	want := "12345"

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestGetPort_WithoutEnvVar(t *testing.T) {
	os.Unsetenv("PORT")

	got := GetPort()
	want := "9000"

	if got != want {
		t.Errorf("expected default %s, got %s", want, got)
	}
}
