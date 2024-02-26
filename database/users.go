package database

import (
	"context"

	"auth/domain/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *MongoUserRepository {
	return &MongoUserRepository{client}
}

type user struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashedPassword"`
	Name           string             `bson:"name"`
}

func mapUser(mongoModel user) (*users.User, error) {
	return &users.User{
		ID:             mongoModel.ID.String(),
		Email:          mongoModel.Email,
		HashedPassword: mongoModel.HashedPassword,
		Name:           mongoModel.Name,
	}, nil
}

func (repo *MongoUserRepository) GetUsersCollection() *mongo.Collection {
	return repo.client.Database("auth").Collection("users")
}

func (repo *MongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	collection := repo.GetUsersCollection()
	var result user
	err := collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return mapUser(result)
}

func (repo *MongoUserRepository) InsertUser(ctx context.Context, input users.ValidCreateUser) (*users.User, error) {
	collection := repo.GetUsersCollection()
	user := user{
		ID:             primitive.ObjectID{},
		Email:          input.Email,
		HashedPassword: input.HashedPassword,
		Name:           input.Name,
	}
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		panic("No inserted id for new user")
	}

	user.ID = insertedID
	return mapUser(user)
}
