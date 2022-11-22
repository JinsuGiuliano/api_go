package store

import (
	"example/api_go/model"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Store) GetValidAccessToken(token string, checkExpiry bool) (qToken *model.AccessToken) {
	tokensCollection := s.db.Collection("tokens")

	var userToken *model.AccessToken

	filter := bson.D{{Key: "token", Value: token}}

	err := tokensCollection.FindOne(s.ctx, filter).Decode(&userToken)
	if err != nil {
		fmt.Println(err.Error())
		return &model.AccessToken{}
	}

	return userToken
}

func (s *Store) CreateAccessToken(token *model.AccessToken) (result *model.AccessToken) {
	tokensCollection := s.db.Collection("tokens")

	_, err := tokensCollection.InsertOne(s.ctx, token)
	if err != nil {
		fmt.Println("token NO Creado")
		panic(err)
	}
	return token
}

func (s *Store) AccessTokens() (t []*model.AccessToken) {
	tokensCollection := s.db.Collection("tokens")

	var tokens []*model.AccessToken

	cursor, err := tokensCollection.Find(s.ctx, bson.M{})
	if err != nil {
		fmt.Println(err.Error())
	}

	if err = cursor.All(s.ctx, &tokens); err != nil {
		fmt.Println(err.Error())
	}

	return tokens
}

func (s *Store) SetAccessTokenToUser(token *model.AccessToken) error {
	usersCollection := s.db.Collection("users")

	filter := bson.D{{Key: "email", Value: token.Email}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: token.Token}}}}

	_, err := usersCollection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteAccessToken(email string) error {
	usersCollection := s.db.Collection("tokens")

	filter := bson.D{{Key: "email", Value: email}}

	result, err := usersCollection.DeleteOne(s.ctx, filter)
	if err != nil {
		fmt.Println("Erorr while deleting token")
		return err
	}
	if result.DeletedCount == 0 {
		fmt.Println("there where no token with the specified email")
		return nil
	}
	return nil
}
