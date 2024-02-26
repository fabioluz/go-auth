package database

import (
	"auth/domain/logs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLogRepository struct {
	client *mongo.Client
}

func NewLogRepository(client *mongo.Client) *MongoLogRepository {
	return &MongoLogRepository{client}
}

type log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	Operation string             `bson:"operation"`
	Timestamp time.Time          `bson:"timestamp"`
}

func mapLog(mongoModel log) *logs.Log {
	return &logs.Log{
		ID:        mongoModel.ID.String(),
		UserID:    mongoModel.UserID,
		Operation: logs.LogOperation(mongoModel.Operation),
		Timestamp: mongoModel.Timestamp,
	}
}

func (repo *MongoLogRepository) GetLogsCollection() *mongo.Collection {
	return repo.client.Database("auth").Collection("logs")
}

func (repo *MongoLogRepository) InsertLog(ctx context.Context, userID string, op logs.LogOperation) (*logs.Log, error) {
	collection := repo.GetLogsCollection()
	log := log{
		ID:        primitive.ObjectID{},
		UserID:    userID,
		Operation: string(op),
		Timestamp: time.Now().UTC(),
	}
	result, err := collection.InsertOne(ctx, log)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		panic("No inserted id for new user")
	}

	log.ID = insertedID
	return mapLog(log), nil
}
