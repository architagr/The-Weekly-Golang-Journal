package model

// Node represents a member in a distributed cluster.
type Node struct {
	ID      string `json:"id"`      // unique identifier
	Address string `json:"address"` // cluster communication address (IP/hostname)
	Port    int    `json:"port"`    // service port
}
