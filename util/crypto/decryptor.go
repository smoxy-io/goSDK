package crypto

type Decrypter interface {
	Decrypt(data []byte) ([]byte, error)
	DecryptString(data string) (string, error)
}
