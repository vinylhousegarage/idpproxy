package store

type Ciphertext struct {
	KID  string `firestore:"kid"`
	Blob string `firestore:"blob"`
}

type TokenEncryptor interface {
	EncryptString(plain string) (Ciphertext, error)
	DecryptString(ct Ciphertext) (string, error)
}
