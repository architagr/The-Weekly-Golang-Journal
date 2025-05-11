package persistence

import (
	"context"
	"error-propagation/pkg/entities"
	"fmt"
	"runtime/debug"
	"time"
)

type UserPersistence struct {
	baseDb *BaseDB
}

var UserPersistenceObj UserPersistence

func init() {
	UserPersistenceObj = UserPersistence{
		baseDb: &baseDbObj,
	}
}

func (usrPer *UserPersistence) Get(ctx context.Context, id int) (*entities.User, error) {
	sqlQuery := fmt.Sprintf("select * from users where id = %d", id)
	data, err := usrPer.baseDb.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, &ObjectNotFoundError{
			ObjectType:        "User",
			ObjectIdentifiers: map[string]any{"id": id},
			StackTrace:        string(debug.Stack()),
		}
	}
	fmt.Println(data)
	// else we will extract data from the data, for this we will hard code this return data
	return &entities.User{
		Id:          id,
		Name:        fmt.Sprintf("user-%d", id),
		CreatedDate: time.Now().Add(-10 * time.Hour),
	}, nil
}
