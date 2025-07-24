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
			log.Printf("Spreading gossip message to peers: %s, message %s", url, payload)
			http.Post(url, "application/json", bytes.NewBuffer(payload))
		}(peer + "/gossip")
	}
}

type PullSpreadStrategy struct{}

func (p *PullSpreadStrategy) Spread(_ model.GossipMessage, peers []string) {
	for _, peer := range peers {
		go func(url string) {
			res, err := http.Get(url)
			if err != nil {
				log.Println("Failed to pull from", url)
				return
			}
			defer res.Body.Close()
			io.ReadAll(res.Body)
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
