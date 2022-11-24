package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUser interface {
	GetUsers() []User
	InsertUser(User) User
	DeleteUser(string) string
	UserBy(string, string) (User, error)
	GetUserByToken(string) (User, error)
	GetUserByID(string) (User, error)
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username,omitempty"`
	Password       string             `bson:"password,omitempty"`
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
