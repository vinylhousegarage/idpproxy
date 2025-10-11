package store

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapNotFound(err error) error {
	if status.Code(err) == codes.NotFound {
		return ErrNotFound
	}

	return err
}

func mapConflict(err error) error {
	switch status.Code(err) {
	case codes.AlreadyExists,
		codes.FailedPrecondition,
		codes.Aborted:
		return ErrConflict
	default:
		return err
	}
}
