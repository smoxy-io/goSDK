package crypto

type SymmetricKey interface {
	Encrypter
	Decrypter

	Key() []byte
	KeyString() string
}

func Encrypt(key Encrypter, data []byte) ([]byte, error) {
	return key.Encrypt(data)
}

func EncryptString(key Encrypter, data string) (string, error) {
	return key.EncryptString(data)
}

func Decrypt(key Decrypter, data []byte) ([]byte, error) {
	return key.Decrypt(data)
}

func DecryptString(key Decrypter, data string) (string, error) {
	return key.DecryptString(data)
}
