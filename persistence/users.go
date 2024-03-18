package persistence

import (
	"context"
	"errors"

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
		ID:             mongoModel.ID.Hex(),
		Email:          mongoModel.Email,
		HashedPassword: mongoModel.HashedPassword,
		Name:           mongoModel.Name,
	}, nil
}

func (repo *MongoUserRepository) Collection() *mongo.Collection {
	return repo.client.Database("auth").Collection("users")
}

func (repo *MongoUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result user
	err = repo.Collection().FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return mapUser(result)
}

func (repo *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	collection := repo.Collection()
	var result user
	err := collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return mapUser(result)
}

func (repo *MongoUserRepository) Insert(ctx context.Context, input users.ValidCreateUser) (*users.User, error) {
	collection := repo.Collection()
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

func (repo *MongoUserRepository) Update(ctx context.Context, id string, input users.ValidUpdateUser) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateFields := bson.M{"$set": bson.M{"name": input.Name}}
	result, err := repo.Collection().UpdateByID(ctx, objectID, updateFields)
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return errors.New("could not update user")
	}

	return nil
}
