package service

import "github.com/vinylhousegarage/idpproxy/internal/authcode/store"

type Service struct {
	store store.Store
}
