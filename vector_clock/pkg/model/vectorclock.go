package model

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// VectorClock tracks causality between distributed events for conflict resolution.
type VectorClock map[string]int

// Copy creates a deep copy of the VectorClock.
// Useful to avoid unintended modifications when passing clocks.
func (vc VectorClock) Copy() VectorClock {
	n := make(VectorClock)
	for k, v := range vc {
		n[k] = v
	}
	return n
}

// Increment increases the counter for the given nodeID by one.
// Should be called on local update events before storing/replicating.
func (vc VectorClock) Increment(nodeID string) {
	vc[nodeID]++
	// Example debug log — enable/disable as needed.
	log.Printf("[DEBUG] VectorClock incremented: node=%s newClock=%s", nodeID, vc.String())
}

// Merge merges the receiver and other vector clocks, returning a new merged clock.
// Merging ensures causal history packages together the latest knowledge from both.
func (vc VectorClock) Merge(other VectorClock) VectorClock {
	out := vc.Copy()
	for id, v := range other {
		if out[id] < v {
			out[id] = v
		}
	}
	log.Printf("[DEBUG] VectorClock merged: base=%s other=%s result=%s", vc.String(), other.String(), out.String())
	return out
}

// Compare compares two vector clocks and returns:
// -1 if receiver vc happened before other,
//
//	1 if vc happened after other,
//	0 if equal,
//	2 if clocks are concurrent (conflict).
func (vc VectorClock) Compare(other VectorClock) int {
	less, greater := false, false
	allIDs := make(map[string]struct{})
	for id := range vc {
		allIDs[id] = struct{}{}
	}
	for id := range other {
		allIDs[id] = struct{}{}
	}
	for id := range allIDs {
		a := vc[id]
		b := other[id]
		if a < b {
			less = true
		}
		if a > b {
			greater = true
		}
	}
	switch {
	case less && !greater:
		return -1
	case greater && !less:
		return 1
	case !less && !greater:
		return 0
	default:
		return 2
	}
}

// String returns a human-readable sorted string of the vector clock.
// e.g. "{node1:3 node2:1}"
func (vc VectorClock) String() string {
	if len(vc) == 0 {
		return "{}"
	}
	var buf []string
	keys := make([]string, 0, len(vc))
	for node := range vc {
		keys = append(keys, node)
	}
	sort.Strings(keys)
	for _, node := range keys {
		buf = append(buf, fmt.Sprintf("%s:%d", node, vc[node]))
	}
	return "{" + strings.Join(buf, " ") + "}"
}
