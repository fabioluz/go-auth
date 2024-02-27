package logs

import (
	"context"
	"time"
)

type LogOperation string

const (
	LogOperationCreateUser   LogOperation = "create-user"
	LogOperationUpdateUser   LogOperation = "update-user"
	LogOperationDeleteDelete LogOperation = "delete-user"
)

type Log struct {
	ID        string       `json:"id"`
	UserID    string       `json:"-"`
	Operation LogOperation `json:"operation"`
	Timestamp time.Time    `json:"timestamp"`
}

type LogRepository interface {
	InsertLog(ctx context.Context, userID string, op LogOperation) (*Log, error)
}
