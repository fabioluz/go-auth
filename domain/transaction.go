package domain

import "context"

type TransactionRepository interface {
	WithTransaction(fn func(ctx context.Context) error) error
}
