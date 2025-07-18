package model

import "time"

type GoogleToken struct {
	AccessToken  string    `firestore:"access_token"`
	RefreshToken string    `firestore:"refresh_token"`
	CreatedAt    time.Time `firestore:"created_at"`
}
