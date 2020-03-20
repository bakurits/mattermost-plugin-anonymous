package main

import (
	"errors"
)

//store public key in plugin's KeyValue Store
func (p *Plugin) storePublicKey(publicKey []byte, userID string) error {
	if err := p.API.KVSet(userID, publicKey); err != nil {
		return errors.New(err.Message)
	}
	return nil
}

//get public key from plugin's KeyValue Store
func (p *Plugin) getPublicKey(userID string) ([]byte, error) {
	userBytes, err := p.API.KVGet(userID)
	if err != nil {
		return nil, errors.New(err.Message)
	}
	publicKey := userBytes
	if publicKey == nil || len(publicKey) == 0 {
		return nil, errors.New("public key is empty")
	}
	return publicKey, nil
}
