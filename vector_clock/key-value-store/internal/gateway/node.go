package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"vectory_clock/pkg/model"
)

var (
	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = fmt.Errorf("resource not found")
)

// Node is a remote node for HTTP-based reads/writes.
type Node struct {
	identifier  string
	fullAddress *url.URL // e.g., http://127.0.0.1:8081
}

// NewNode constructs a node from id, address, and port.
func NewNode(identifier, address string, port int) (*Node, error) {
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}
	return &Node{
		identifier:  identifier,
		fullAddress: url,
	}, nil
}

// GetIdentifier returns node's cluster-unique ID.
func (n *Node) GetIdentifier() string {
	return n.identifier
}

// GetFullAddress returns the (cached) HTTP endpoint for the node.
func (n *Node) GetFullAddress() string {
	return n.fullAddress.String()
}

// GetValue fetches a value (with vector clock) from the remote node via GET.
func (n *Node) GetValue(k string) (*model.ValueWithClock, error) {
	var v *model.ValueWithClock
	url := n.fullAddress.String() + "/" + k
	log.Printf("[CLIENT][%s] GET %s", n.identifier, url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("[CLIENT][%s][ERROR] crafting GET: %v", n.identifier, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CLIENT][%s][ERROR] request GET/%s: %v", n.identifier, k, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[CLIENT][%s][INFO] GET %s: not found", n.identifier, k)
		return nil, ErrNotFound
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("[CLIENT][%s][ERROR] GET %s: server error %v", n.identifier, k, resp.Status)
		return nil, fmt.Errorf("non-200 response: %v", resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Printf("[CLIENT][%s][ERROR] decoding GET %s: %v", n.identifier, k, err)
		return nil, err
	}
	log.Printf("[CLIENT][%s] GET %s: value=%v clock=%v", n.identifier, k, v.Value, v.Clock)
	return v, nil
}

// SetValueWithClock sends a value (with vector clock) to the node using PUT.
func (n *Node) SetValueWithClock(key string, v *model.ValueWithClock) (*model.ValueWithClock, error) {
	body, err := json.Marshal(v)
	if err != nil {
		log.Printf("[CLIENT][%s][ERROR] marshal PUT body: %v", n.identifier, err)
		return nil, err
	}
	url := n.fullAddress.String() + "/" + key
	log.Printf("[CLIENT][%s] PUT %s: value=%v clock=%v", n.identifier, key, v.Value, v.Clock)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[CLIENT][%s][ERROR] crafting PUT: %v", n.identifier, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CLIENT][%s][ERROR] PUT request to %s: %v", n.identifier, url, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[CLIENT][%s][WARN] PUT %s: not found", n.identifier, key)
		return nil, ErrNotFound
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("[CLIENT][%s][ERROR] PUT %s: server error %v", n.identifier, key, resp.Status)
		return nil, fmt.Errorf("non-200 response: %v", resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Printf("[CLIENT][%s][ERROR] decoding PUT response: %v", n.identifier, err)
		return nil, err
	}
	log.Printf("[CLIENT][%s] PUT %s succeeded", n.identifier, key)
	return v, nil
}
