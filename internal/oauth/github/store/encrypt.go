package store

type Ciphertext struct {
	KID  string
	Blob string
}

type TokenEncryptor interface {
	EncryptString(plain string) (Ciphertext, error)
	DecryptString(ct Ciphertext) (string, error)
}
