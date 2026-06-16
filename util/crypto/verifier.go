package crypto

type Verifier interface {
	VerifySignature(data []byte, sig []byte) bool
	VerifySignatureString(data string, sig string) bool
}
