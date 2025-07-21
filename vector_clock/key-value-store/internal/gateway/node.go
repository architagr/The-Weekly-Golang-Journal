package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"vectory_clock/pkg/model"
)

var (
	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = fmt.Errorf("resource not found")
)

type Node struct {
	identifier  string
	fullAddress *url.URL
}

// NewNode creates a new Node with the given identifier, address, and port.
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

// GetIdentifier returns the identifier of the Node.
func (n *Node) GetIdentifier() string {
	return n.identifier
}

func (n *Node) GetFullAddress() string {
	return n.fullAddress.String()
}

func (n *Node) GetValue(k string) (*model.ValueWithClock, error) {
	var v *model.ValueWithClock
	req, err := http.NewRequest(http.MethodGet, n.fullAddress.String()+"/"+k, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func (n *Node) SetValueWithClock(key string, v *model.ValueWithClock) (*model.ValueWithClock, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, n.fullAddress.String()+"/"+key, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
