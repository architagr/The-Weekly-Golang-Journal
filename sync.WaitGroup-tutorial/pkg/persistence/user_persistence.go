package persistence

import (
	"context"
	"fmt"
	"strings"
	"waitGroup-tutorial/pkg/entities"
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

func (usrPer *UserPersistence) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	sqlQuery := fmt.Sprintf("insert into users (name, email, tags) values(%s, %s, %s)", user.Name, user.Email, strings.Join(user.NotificationTags, ","))
	id, err := usrPer.baseDb.Insert(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	user.Id = id
	// else we will extract data from the data, for this we will hard code this return data
	return user, nil
}
