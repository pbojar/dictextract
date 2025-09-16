package dawg

import (
	"fmt"
	"sort"
	"strings"
)

type DAWGNode struct {
	id         int
	isTerminal bool
	children   map[rune]*DAWGNode
}

// signature returns a unique string for a DAWGNode based on its children. Used
// in the DAWGBuilder minimize routine to track nodes registered by the builder.
func (n *DAWGNode) signature() string {
	var sb strings.Builder
	if n.isTerminal {
		sb.WriteString("1_")
	} else {
		sb.WriteString("0_")
	}

	// Sort keys for a deterministic signature
	keys := make([]rune, 0, len(n.children))
	for r := range n.children {
		keys = append(keys, r)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, r := range keys {
		sb.WriteRune(r)
		sb.WriteString(fmt.Sprintf("%d_", n.children[r].id))
	}
	return sb.String()
}
