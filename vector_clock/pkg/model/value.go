package model

type ValueWithClock struct {
	Value any         `json:"value"`
	Clock VectorClock `json:"clock"`
}
