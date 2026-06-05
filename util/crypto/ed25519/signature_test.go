package ed25519

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/ed25519"
)

func Test_Sign(t *testing.T) {
	pubKey, privKey, pkErr := ed25519.GenerateKey(rand.Reader)

	if pkErr != nil {
		t.Errorf("error generating ed25519 key: %v", pkErr)
		return
	}

	data := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus semper accumsan neque id faucibus. Mauris ullamcorper felis vel hendrerit efficitur. Sed eu ante ac augue aliquam rhoncus. Cras et laoreet arcu. Mauris porttitor malesuada dui, eu tincidunt magna mattis non. Mauris quis aliquam lacus. Phasellus ut finibus augue. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Suspendisse facilisis condimentum nunc, id mattis sem ullamcorper pulvinar. Quisque nec magna nisi. Cras ante orci, tempor at diam a, hendrerit vehicula risus."

	b64Data := base64.StdEncoding.EncodeToString([]byte(data))

	sig, sErr := SignString(privKey, b64Data)

	if sErr != nil {
		t.Errorf("error signing data: %v", sErr)
		return
	}

	if !VerifySignatureString(pubKey, b64Data, sig) {
		t.Errorf("failed to verify signature")
		return
	}
}
