package signer

import "time"

const (
	testKey     = "secret"
	testKid123  = "kid-123"
	testKidXYZ  = "kid-xyz"
	testKidErr  = "kid-err"
	testKidErr2 = "kid-err2"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }
