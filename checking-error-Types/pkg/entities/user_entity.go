package entities

import (
	"time"

	list "github.com/architagr/golang_collections/list"
)

type User struct {
	Id          int
	Name        string
	CreatedDate time.Time
	Sessions    list.IItratorList[*Session, Session]
}

type Session struct {
	Id        int
	UserId    int
	StartTime time.Time
	EndTime   time.Time
}

func (sess *Session) Copy() interface{} {
	cpy := new(Session)
	cpy.Id = sess.Id
	cpy.UserId = sess.UserId
	cpy.StartTime = sess.StartTime
	cpy.EndTime = sess.EndTime
	return cpy
}
func (sess *Session) Equal(val interface{}) bool {
	data, ok := val.(*Session)
	if !ok {
		return false
	}
	return data.Id == sess.Id && data.UserId == sess.UserId && data.StartTime.Equal(sess.StartTime) && data.EndTime.Equal(sess.EndTime)
}
