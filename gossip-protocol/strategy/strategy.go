package strategy

import (
	"gossip-protocol/config"
	"gossip-protocol/model"
	"math/rand"
	"time"
)

type GossipStrategy interface {
	GenerateMessage(state map[string]bool) model.GossipMessage
	Merge(local map[string]bool, incoming model.GossipMessage) map[string]bool
}

type AntiEntropyStrategy struct{}

func (a *AntiEntropyStrategy) GenerateMessage(state map[string]bool) model.GossipMessage {
	return model.GossipMessage{
		SenderID:   config.SelfID,
		Timestamp:  time.Now(),
		NodeHealth: state,
		Config:     *config.CurrentConfig,
	}
}
func (a *AntiEntropyStrategy) Merge(local map[string]bool, incoming model.GossipMessage) map[string]bool {
	return incoming.NodeHealth
}

type RumorMongeringStrategy struct{}

func (r *RumorMongeringStrategy) GenerateMessage(state map[string]bool) model.GossipMessage {
	partial := map[string]bool{}
	for k, v := range state {
		if rand.Intn(2) == 0 {
			partial[k] = v
		}
	}
	return model.GossipMessage{
		SenderID:   config.SelfID,
		Timestamp:  time.Now(),
		NodeHealth: partial,
		Config:     *config.CurrentConfig,
	}
}
func (r *RumorMongeringStrategy) Merge(local map[string]bool, incoming model.GossipMessage) map[string]bool {
	for k, v := range incoming.NodeHealth {
		local[k] = v
	}
	return local
}

type AggregationStrategy struct{}

func (a *AggregationStrategy) GenerateMessage(state map[string]bool) model.GossipMessage {
	summary := map[string]bool{"healthy": len(state) > 0}
	return model.GossipMessage{
		SenderID:   config.SelfID,
		Timestamp:  time.Now(),
		NodeHealth: summary,
		Config:     *config.CurrentConfig,
	}
}
func (a *AggregationStrategy) Merge(local map[string]bool, incoming model.GossipMessage) map[string]bool {
	return local
}

func GetGossipStrategy(name string) GossipStrategy {
	switch name {
	case "anti-entropy":
		return &AntiEntropyStrategy{}
	case "rumor-mongering":
		return &RumorMongeringStrategy{}
	case "aggregation":
		return &AggregationStrategy{}
	default:
		return &AntiEntropyStrategy{}
	}
}
