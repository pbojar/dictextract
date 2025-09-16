package dawg

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

// SerializableDAWGNode is a serializable version of DAWGNode.
// Its children map is a map of integer IDs that replaces the map of pointers.
type SerializableDAWGNode struct {
	IsTerminal bool
	Children   map[rune]int
}

// SerializableDAWG is the main structure for serialization.
type SerializableDAWG struct {
	RootID int
	Nodes  []SerializableDAWGNode
}

// SaveAsGob writes the DAWG to a file at the given path using gob encoding.
func (d *DAWG) SaveAsGob(path string) error {
	// Traverse the graph to flatten it into a list of nodes.
	// Visited nodes are tracked in a map and a slice is used as a queue for BFS.
	nodes := []*DAWGNode{}
	nodeMap := make(map[int]*DAWGNode)
	queue := []*DAWGNode{d.root}
	nodeMap[d.root.id] = d.root

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		nodes = append(nodes, current)

		for _, child := range current.children {
			if _, visited := nodeMap[child.id]; !visited {
				nodeMap[child.id] = child
				queue = append(queue, child)
			}
		}
	}

	// Sort nodes by ID to ensure a deterministic order.
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].id < nodes[j].id
	})

	// Create a map from original ID to the new index in the slice.
	idToIndex := make(map[int]int)
	for i, node := range nodes {
		idToIndex[node.id] = i
	}

	// Create the serializable representation.
	sNodes := make([]SerializableDAWGNode, len(nodes))
	for i, node := range nodes {
		sChildren := make(map[rune]int)
		for r, child := range node.children {
			sChildren[r] = idToIndex[child.id] // Use the new index as the ID
		}
		sNodes[i] = SerializableDAWGNode{
			IsTerminal: node.isTerminal,
			Children:   sChildren,
		}
	}

	sDAWG := SerializableDAWG{
		RootID: idToIndex[d.root.id],
		Nodes:  sNodes,
	}

	// Write to file using gob.
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(sDAWG); err != nil {
		return fmt.Errorf("failed to encode DAWG: %w", err)
	}
	return writer.Flush()
}

// LoadDAWGFromGob reads a gob-encoded DAWG from a file and reconstructs it.
func LoadDAWGFromGob(path string) (*DAWG, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)

	var sDAWG SerializableDAWG
	if err := decoder.Decode(&sDAWG); err != nil {
		return nil, fmt.Errorf("failed to decode DAWG: %w", err)
	}

	// Reconstruct DAWGNodes from the serialized nodes.
	// First pass: create all node objects and store them in a slice.
	nodes := make([]*DAWGNode, len(sDAWG.Nodes))
	for i := range sDAWG.Nodes {
		nodes[i] = &DAWGNode{
			id:         i, // The new ID is the slice index
			isTerminal: sDAWG.Nodes[i].IsTerminal,
			children:   make(map[rune]*DAWGNode),
		}
	}

	// Second pass: reconstruct the children pointers.
	for i, sNode := range sDAWG.Nodes {
		for r, childIndex := range sNode.Children {
			nodes[i].children[r] = nodes[childIndex]
		}
	}

	// Create the final DAWG object.
	dawg := &DAWG{
		root: nodes[sDAWG.RootID],
	}

	return dawg, nil
}
