package session

import "time"

type Usecase struct {
	Repo Repository
	Now  func() time.Time
}
