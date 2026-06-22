package gcm256

import (
	"testing"
)

func Test_EncryptDecrypt(t *testing.T) {
	privKey := NewKey(nil)

	if !privKey.IsValid() {
		t.Errorf("failed to generate private key")
		return
	}

	data := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus semper accumsan neque id faucibus. Mauris ullamcorper felis vel hendrerit efficitur. Sed eu ante ac augue aliquam rhoncus. Cras et laoreet arcu. Mauris porttitor malesuada dui, eu tincidunt magna mattis non. Mauris quis aliquam lacus. Phasellus ut finibus augue. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Suspendisse facilisis condimentum nunc, id mattis sem ullamcorper pulvinar. Quisque nec magna nisi. Cras ante orci, tempor at diam a, hendrerit vehicula risus."

	encData, edErr := EncryptString(privKey.KeyString(), data)

	if edErr != nil {
		t.Errorf("error encrypting data: %v", edErr)
		return
	}

	decData, ddErr := DecryptString(privKey.KeyString(), encData)

	if ddErr != nil {
		t.Errorf("error decrypting data: %v", ddErr)
		return
	}

	if decData != data {
		t.Errorf("decrypted data does not match original data\nexpected:\n%s\n\ngot: %s", data, decData)
		return
	}
}
