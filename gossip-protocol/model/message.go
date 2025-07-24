package model

import (
	"gossip-protocol/config"
	"time"

	"github.com/google/uuid"
)

type GossipMessage struct {
	SenderID   uuid.UUID           `json:"senderId"`
	Timestamp  time.Time           `json:"timestamp"`
	NodeHealth map[string]bool     `json:"nodeHealth"`
	Config     config.GossipConfig `json:"config"`
}
