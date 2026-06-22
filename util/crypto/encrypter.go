package crypto

type Encrypter interface {
	Encrypt(data []byte) ([]byte, error)
	EncryptString(data string) (string, error)
}
