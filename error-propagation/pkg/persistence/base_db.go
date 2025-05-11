package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrQueryCancelled = errors.New("query cancelled")
)
var baseDbObj BaseDB

func init() {
	baseDbObj = *NewBaseDB(2 * time.Second)
}

// BaseDB represents a mock database client.
type BaseDB struct {
	queryDelay time.Duration
}

// NewBaseDB creates a new BaseDB instance.
func NewBaseDB(delay time.Duration) *BaseDB {
	return &BaseDB{queryDelay: delay}
}

// Query simulates a database query that respects context cancellation.
func (db *BaseDB) Query(ctx context.Context, query string) (string, error) {
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		// Simulate a query delay
		// Ideally here we will create a connection and then executes the query
		time.Sleep(db.queryDelay)
		if ctx.Err() != nil {
			errChan <- fmt.Errorf("%w:%s", ErrQueryCancelled, ctx.Err().Error())
			return
		}
		fmt.Println("query executed")
		// Simulate query result
		resultChan <- fmt.Sprintf("Result for query: %s", query)
	}()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("%w:%s", ErrQueryCancelled, ctx.Err().Error())
	case err := <-errChan:
		return "", err
	case result := <-resultChan:
		return result, nil
	}
}
