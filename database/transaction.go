package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTransactionRepository struct {
	client *mongo.Client
}

func NewTransactionRepository(client *mongo.Client) *MongoTransactionRepository {
	return &MongoTransactionRepository{client}
}

func (repo *MongoTransactionRepository) WithTransaction(fn func(ctx context.Context) error) error {
	return repo.client.UseSession(context.TODO(), func(sessionCtx mongo.SessionContext) error {
		_, err := sessionCtx.WithTransaction(context.TODO(), func(tranCtx mongo.SessionContext) (interface{}, error) {
			err := fn(tranCtx)
			return nil, err
		})
		return err
	})
}
