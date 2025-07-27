package config

import (
	"time"

	"github.com/google/uuid"
)

var (
	SelfID       uuid.UUID
	Port         string
	Strategy     string
	SpreadMethod string
	Peers        []string
)

func init() {
	// Generate unique node ID
	SelfID, _ = uuid.NewV7()
	Strategy = "anti-entropy"
	SpreadMethod = "push"
	Peers = []string{}
}

type GossipConfig struct {
	Fanout     int           `json:"fanout"`
	Interval   time.Duration `json:"interval"`
	BufferSize int           `json:"bufferSize"`
}

var CurrentConfig *GossipConfig

func DefaultConfig() *GossipConfig {
	CurrentConfig = &GossipConfig{
		Fanout:     3,
		Interval:   5 * time.Second,
		BufferSize: 2,
	}
	return CurrentConfig
}
