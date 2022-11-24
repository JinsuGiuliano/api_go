package model

import (
	"time"
)

type AccessToken struct {
	Token  string    `bson:"token"`
	Id     string    `bson:"id"`
	Expire time.Time `bson:"expire"`
	Email  string    `bson:"email"`
}

type IAccessTokenData interface {
	AccessTokens() []*AccessToken
	CreateAccessToken(*AccessToken) *AccessToken
	GetValidAccessToken(token string, checkExpiry bool) *AccessToken
	SetAccessTokenToUser(*AccessToken) error
	DeleteAccessToken(string) error
}
