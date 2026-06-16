package crypto

type Key interface {
	Signer
	Verifier
}

func Sign(key Signer, data []byte) ([]byte, error) {
	return key.Sign(data)
}

func SignString(key Signer, data string) (string, error) {
	return key.SignString(data)
}

func VerifySignature(key Verifier, data []byte, sig []byte) bool {
	return key.VerifySignature(data, sig)
}

func VerifySignatureString(key Verifier, data string, sig string) bool {
	return key.VerifySignatureString(data, sig)
}
