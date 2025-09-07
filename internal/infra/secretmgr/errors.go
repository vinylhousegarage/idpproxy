package secretmgr

import "errors"

var (
	ErrAccessFailed = errors.New("secretmgr: access failed")
	ErrInitFailed   = errors.New("secretmgr: init failed")
)
