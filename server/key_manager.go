package main

import (
	"errors"
)

//store public key in plugin's KeyValue Store
func (p *Plugin) storePublicKey(publicKey string, userId string) error {
	pb := []byte(publicKey)
	if err := p.API.KVSet(userId, pb); err != nil {
		return errors.New(err.Message)
	}
	return nil
}

//get public key from plugin's KeyValue Store
func (p *Plugin) getPublicKey(userId string) (string, error) {
	userBytes, err := p.API.KVGet(userId)
	if err != nil {
		return "", errors.New(err.Message)
	}
	publicKey := string(userBytes)
	if publicKey == "" {
		return "", errors.New("public key is empty")
	}
	return publicKey, nil
}
