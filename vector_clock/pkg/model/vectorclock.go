package model

import (
	"fmt"
	"sort"
	"strings"
)

type VectorClock map[string]int

func (vc VectorClock) Copy() VectorClock {
	n := make(VectorClock)
	for k, v := range vc {
		n[k] = v
	}
	return n
}

func (vc VectorClock) Increment(nodeID string) {
	vc[nodeID]++
}

func (vc VectorClock) Merge(other VectorClock) VectorClock {
	out := vc.Copy()
	for id, v := range other {
		if out[id] < v {
			out[id] = v
		}
	}
	return out
}

// Returns:
// -1 if vc < other (happened-before)
// 1 if vc > other (happened-after)
// 0 if equal
// 2 if concurrent/conflicting
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
