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
