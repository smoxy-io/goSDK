package crypto

type Signer interface {
	Sign(data []byte) ([]byte, error)
	SignString(data string) (string, error)
}
