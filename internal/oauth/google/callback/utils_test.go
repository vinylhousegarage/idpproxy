package callback

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateSessionID(t *testing.T) {
	id := generateSessionID()

	require.NotEmpty(t, id, "UUID should not be empty")

	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[a-f0-9]{12}$`)
	require.True(t, uuidRegex.MatchString(id), "UUID format is invalid")
}
