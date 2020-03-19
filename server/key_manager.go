package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mattermost/mattermost-server/v5/model"
)

type AnonymousUser struct {
	model.User
	PublicKey string `json:"public_key"`
}

//store public key in plugin's KeyValue Store
func (p *Plugin) storePublicKey(publicKey string, userId string) error {
	modelUser, err := p.API.GetUser(userId)
	if err != nil || modelUser == nil {
		return err
	}

	var user = AnonymousUser{
		User:      *modelUser,
		PublicKey: publicKey,
	}

	reqBodyBytes := new(bytes.Buffer)
	errno := json.NewEncoder(reqBodyBytes).Encode(user)
	if errno != nil {
		return err
	}

	return p.API.KVSet(userId, reqBodyBytes.Bytes())
}

//get public key from plugin's KeyValue Store
func (p *Plugin) getPublicKey(userId string) (string, error) {
	userBytes, err := p.API.KVGet(userId)
	if err != nil {
		return "", err
	}

	var user AnonymousUser
	errno := json.NewDecoder(bytes.NewReader(userBytes)).Decode(&user)
	if errno != nil {
		return "", err
	}

	if user.PublicKey == "" {
		return "", errors.New("public key is empty")
	}

	return user.PublicKey, nil
}
