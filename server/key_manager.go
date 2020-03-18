package main

import (
	"bytes"
	"encoding/json"
	"github.com/mattermost/mattermost-server/v5/model"
)

type AnonymousUser struct {
	model.User
	PublicKey string `json:"public_key"`
}

//store public key in plugin's KeyValue Store
func (p *Plugin) storePublicKey(publicKey string, user_id string) *model.AppError {
	model_user, err := p.API.GetUser(user_id)
	if err != nil || model_user == nil {
		return err
	}
	var user = AnonymousUser{
		User:      *model_user,
		PublicKey: publicKey,
	}

	reqBodyBytes := new(bytes.Buffer)
	_ = json.NewEncoder(reqBodyBytes).Encode(user)

	return p.API.KVSet(user_id, reqBodyBytes.Bytes())
}

//get public key from plugin's KeyValue Store
func (p *Plugin) getPublicKey(user_id string) (string, *model.AppError) {
	user_bytes, err := p.API.KVGet(user_id)
	var user AnonymousUser
	_ = json.NewDecoder(bytes.NewReader(user_bytes)).Decode(&user)

	return user.PublicKey, err
}
