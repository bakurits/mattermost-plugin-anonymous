package crypto

import (
	"encoding/base64"

	"github.com/pkg/errors"
)

// PublicKey stores public key data
type PublicKey []byte

// PublicKeyFromString decodes string end returns public key
func PublicKeyFromString(key string) (PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return PublicKey{}, errors.Wrap(err, "can't decode public key from data")
	}
	return data, nil
}

// String encodes public key data in base64 formatting
func (pb PublicKey) String() string {
	return base64.StdEncoding.EncodeToString(pb)
}
