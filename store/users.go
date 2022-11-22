package store

import (
	model "example/api_go/model"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) GetUsers() []model.User {
	userCollection := s.db.Collection("users")

	var Users []model.User
	cursor, err := userCollection.Find(s.ctx, bson.M{})
	if err != nil {
		fmt.Println("UserCollection.Find: ", err.Error())
	}

	if err = cursor.All(s.ctx, &Users); err != nil {
		fmt.Println("GetUsers.cursor.All: ", err.Error())
	}

	return Users
}

func (s *Store) GetUserByID(ID string) (u model.User, e error) {
	userCollection := s.db.Collection("users")

	var User model.User

	filter := bson.D{{Key: "_id", Value: ID}}

	contactCursor := userCollection.FindOne(s.ctx, filter)

	err := contactCursor.Err()
	if err != nil {
		fmt.Println("FindOne: ", err.Error())
	}

	return User, e
}

func (s *Store) InsertUser(u model.User) model.User {
	userCollection := s.db.Collection("users")

	_, err := userCollection.InsertOne(s.ctx, u)
	if err != nil {
		fmt.Println("InsertOne: ", err.Error())
		return model.User{}
	}

	return u
}

func (s *Store) DeleteUser(ID string) string {
	userCollection := s.db.Collection("users")

	idPrimitive, ok := primitive.ObjectIDFromHex(ID)
	if ok != nil {
		fmt.Println("ObjectIDFromHex: ", ok.Error())
	}
	_, err := userCollection.DeleteOne(s.ctx, bson.M{"_id": idPrimitive})
	if err != nil {
		fmt.Println("DeleteOne: ", err.Error())
	}
	return ID
}

func (s *Store) GetUserByToken(token string) (model.User, error) {
	var user model.User

	usersCollection := s.db.Collection("users")

	filter := bson.D{{Key: "token", Value: token}}
	err := usersCollection.FindOne(s.ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println("GetUserByToken: ", err.Error())
		return model.User{}, err
	}

	return user, nil
}

func (s *Store) GetUserByGoogleID(id string) (model.User, error) {
	var user model.User

	usersCollection := s.db.Collection("users")

	filter := bson.D{{Key: "GoogleID", Value: id}}
	err := usersCollection.FindOne(s.ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err.Error())
	}

	return user, nil
}

func (s *Store) UserByEmail(email string) (u model.User, e error) {
	var User model.User
	fmt.Println("UserByEmail")
	userCollection := s.db.Collection("users")

	filter := bson.D{{Key: "email", Value: email}}
	err := userCollection.FindOne(s.ctx, filter).Decode(&User)
	if err != nil {
		return model.User{}, fmt.Errorf("there is no user associated with the given email")
	}
	fmt.Println("return User, e" + User.Email)
	return User, e
}
