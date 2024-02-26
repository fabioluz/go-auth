package logs

import (
	"context"
	"time"
)

type LogOperation string

const (
	LogOperationCreate LogOperation = "create"
	LogOperationUpdate LogOperation = "update"
	LogOperationDelete LogOperation = "delete"
)

type Log struct {
	ID        string
	UserID    string
	Operation LogOperation
	Timestamp time.Time
}

type LogRepository interface {
	InsertLog(ctx context.Context, userID string, op LogOperation) (*Log, error)
}
