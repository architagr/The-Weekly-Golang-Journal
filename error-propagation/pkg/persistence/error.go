package persistence

import (
	"encoding/json"
	"fmt"
)

type ObjectNotFoundError struct {
	ObjectType        string         `json:"objectType"`
	ObjectIdentifiers map[string]any `json:"objectIdentifier"`
	StackTrace        string         `json:"stackTrace"`
}

func (e *ObjectNotFoundError) Error() string {
	return fmt.Sprintf("%s not found, filter properties: [%v]", e.ObjectType, e.ObjectIdentifiers)
}
func (e *ObjectNotFoundError) LogMessage() string {
	b, _ := json.Marshal(e)
	return string(b)
}
