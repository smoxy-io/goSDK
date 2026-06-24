package shamir

import (
	"encoding/base64"

	"github.com/hashicorp/vault/shamir"
)

func CreateShamirSecret(secret string, parts int, threshold int) ([]string, error) {
	shares, sErr := shamir.Split([]byte(secret), parts, threshold)

	if sErr != nil {
		return nil, sErr
	}

	s := make([]string, 0, len(shares))

	for _, b := range shares {
		s = append(s, base64.StdEncoding.EncodeToString(b))
	}

	return s, nil
}

func GetShamirSecret(shares []string) (string, error) {
	if len(shares) == 0 {
		return "", errors.New("no shares provided")
	}

	parts := make([][]byte, 0, len(shares))

	for _, share := range shares {
		p, pErr := base64.StdEncoding.DecodeString(share)

		if pErr != nil {
			return "", errors.New("error parsing share: %v", pErr)
		}

		parts = append(parts, p)
	}

	secret, sErr := shamir.Combine(parts)

	if sErr != nil {
		return "", sErr
	}

	return string(secret), nil
}
