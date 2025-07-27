package spread

import (
	"bytes"
	"encoding/json"
	"gossip-protocol/model"
	"io"
	"log"
	"net/http"
)

type SpreadStrategy interface {
	Spread(msg model.GossipMessage, peers []string)
}

type PushSpreadStrategy struct{}

func (p *PushSpreadStrategy) Spread(msg model.GossipMessage, peers []string) {
	for _, peer := range peers {
		go func(url string) {
			payload, _ := json.Marshal(msg)
			log.Printf("[SEND] Sending gossip → %s | Payload: %s", url, string(payload))

			resp, err := http.Post(url+"/gossip", "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("[ERROR] Failed to send gossip to %s: %v", url, err)
				return
			}
			defer resp.Body.Close()
			log.Printf("[ACK] Gossip sent successfully to %s", url)
		}(peer)
	}
}

type PullSpreadStrategy struct{}

func (p *PullSpreadStrategy) Spread(_ model.GossipMessage, peers []string) {
	for _, peer := range peers {
		go func(url string) {
			resp, err := http.Get(url + "/health")
			if err != nil {
				log.Printf("[ERROR] Failed to pull gossip from %s: %v", url, err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			log.Printf("[PULL] Pulled gossip from %s → %s", url, string(body))
		}(peer)
	}
}

type PushPullSpreadStrategy struct{}

func (p *PushPullSpreadStrategy) Spread(msg model.GossipMessage, peers []string) {
	(&PushSpreadStrategy{}).Spread(msg, peers)
	(&PullSpreadStrategy{}).Spread(msg, peers)
}

func GetSpreadStrategy(name string) SpreadStrategy {
	switch name {
	case "push":
		return &PushSpreadStrategy{}
	case "pull":
		return &PullSpreadStrategy{}
	case "push-pull":
		return &PushPullSpreadStrategy{}
	default:
		return &PushSpreadStrategy{}
	}
}
