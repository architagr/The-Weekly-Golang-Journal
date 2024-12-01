package apperrors

import "fmt"

type NotFound struct {
	Id int
}

func (e NotFound) Error() string {
	return fmt.Sprintf("data not found, Id (%d)", e.Id)
}
