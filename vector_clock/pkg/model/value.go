package model

// ValueWithClock encapsulates a stored value with its version metadata - a VectorClock.
// The 'Value' is stored as interface{} (Go 1.18+ 'any') for flexibility in demos.
type ValueWithClock struct {
	Value any         `json:"value"`
	Clock VectorClock `json:"clock"`
}
