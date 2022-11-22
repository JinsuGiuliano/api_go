package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUser interface {
	GetUsers() []User
	InsertUser(User) User
	DeleteUser(string) string
	UserByEmail(string) (User, error)
	GetUserByToken(string) (User, error)
	GetUserByID(string) (User, error)
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty"`
	Email          string             `bson:"email,omitempty"`
	Token          string             `bson:"token,omitempty"`
	AllowTokenAuth bool               `bson:"allowTokenAuth"`
	DataAllowed    bool               `bson:"dataAllowed"`
}

type IGoogleUSer struct {
	Id             string
	Email          string
	Verified_email bool
	Picture        string
}

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
