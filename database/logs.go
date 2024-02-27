package database

import (
	"auth/domain/logs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLogRepository struct {
	client *mongo.Client
}

func NewLogRepository(client *mongo.Client) *MongoLogRepository {
	return &MongoLogRepository{client}
}

type log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId"`
	Operation string             `bson:"operation"`
	Timestamp time.Time          `bson:"timestamp"`
}

func mapLog(mongoModel log) logs.Log {
	return logs.Log{
		ID:        mongoModel.ID.Hex(),
		UserID:    mongoModel.UserID.Hex(),
		Operation: logs.LogOperation(mongoModel.Operation),
		Timestamp: mongoModel.Timestamp,
	}
}

func (repo *MongoLogRepository) GetLogsCollection() *mongo.Collection {
	return repo.client.Database("auth").Collection("logs")
}

func (repo *MongoLogRepository) GetLogs(ctx context.Context, userID string, pageSize int, after string) ([]logs.Log, error) {
	objUserId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	query := bson.M{"userId": objUserId}
	if after != "" {
		objAfter, err := primitive.ObjectIDFromHex(after)
		if err != nil {
			return nil, err
		}
		query["_id"] = bson.M{"$gt": objAfter}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))

	cursor, err := repo.GetLogsCollection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mongoLogs := []log{}
	if err := cursor.All(ctx, &mongoLogs); err != nil {
		return nil, err
	}

	logs := []logs.Log{}
	for _, log := range mongoLogs {
		logs = append(logs, mapLog(log))
	}

	return logs, nil
}

func (repo *MongoLogRepository) InsertLog(ctx context.Context, userID string, op logs.LogOperation) (*logs.Log, error) {
	userObjId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	collection := repo.GetLogsCollection()
	mongoLog := log{
		ID:        primitive.ObjectID{},
		UserID:    userObjId,
		Operation: string(op),
		Timestamp: time.Now().UTC(),
	}
	result, err := collection.InsertOne(ctx, mongoLog)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		panic("No inserted id for new user")
	}

	mongoLog.ID = insertedID

	log := mapLog(mongoLog)
	return &log, nil
}
