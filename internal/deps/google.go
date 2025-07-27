package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/verify"
)

type GoogleDependencies struct {
	Logger   *zap.Logger
	Verifier verify.Verifier
}

func NewGoogleDeps(verifier verify.Verifier, logger *zap.Logger) *GoogleDependencies {
	return &GoogleDependencies{
		Logger:   logger,
		Verifier: verifier,
	}
}
