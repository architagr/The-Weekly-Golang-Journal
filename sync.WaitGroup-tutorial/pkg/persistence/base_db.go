package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"
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

// Query simulates a database insert that respects context cancellation.
func (db *BaseDB) Insert(ctx context.Context, query string) (int, error) {
	resultChan := make(chan int, 1)
	errChan := make(chan error, 1)

	go func() {
		// Simulate a insert delay
		// Ideally here we will create a connection and then executes the query
		time.Sleep(db.queryDelay)
		if ctx.Err() != nil {
			errChan <- ctx.Err()
			return
		}
		fmt.Println("insert executed")
		// Simulate insert result
		resultChan <- 1
	}()

	select {
	case <-ctx.Done():
		return 0, errors.New("query cancelled: " + ctx.Err().Error())
	case err := <-errChan:
		return 0, err
	case result := <-resultChan:
		return result, nil
	}
}
