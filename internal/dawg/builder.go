package dawg

import "fmt"

// DAWGBuilder is used to construct a new DAWG.
// Words MUST be inserted in lexicographical order.
type DAWGBuilder struct {
	root            *DAWGNode
	registeredNodes map[string]*DAWGNode
	lastWord        string
	nodeCounter     int
}

// NewDAWGBuilder creates a new DAWGBuilder.
func NewDAWGBuilder() *DAWGBuilder {
	root := &DAWGNode{
		id:       0,
		children: make(map[rune]*DAWGNode),
	}
	return &DAWGBuilder{
		root:            root,
		registeredNodes: make(map[string]*DAWGNode),
		nodeCounter:     1,
	}
}

// newNode creates a new DAWGNode with a unique ID.
func (b *DAWGBuilder) newNode() *DAWGNode {
	node := &DAWGNode{
		id:       b.nodeCounter,
		children: make(map[rune]*DAWGNode),
	}
	b.nodeCounter++
	return node
}

// Insert adds a word to the DAWG. Words MUST be inserted in
// lexicographical order for the algorithm to work correctly.
func (b *DAWGBuilder) Insert(word string) (err error) {
	if word < b.lastWord {
		return fmt.Errorf("words must be inserted in lexicographical order: received '%s' after '%s'", word, b.lastWord)
	}

	// Find the common prefix length with the last word
	comPreLen := 0
	for comPreLen < len(word) && comPreLen < len(b.lastWord) && word[comPreLen] == b.lastWord[comPreLen] {
		comPreLen++
	}

	// Minimize the suffix of the last word that is not part of the common prefix.
	b.minimize(comPreLen)

	// Add the new suffix for the current word.
	var node *DAWGNode
	if b.lastWord == "" {
		node = b.root
	} else {
		// Find the node where the new suffix should branch off
		node = b.root
		for i := 0; i < comPreLen; i++ {
			node = node.children[rune(b.lastWord[i])]
		}
	}

	// Add the new nodes for the current word's suffix.
	for _, r := range word[comPreLen:] {
		nextNode := b.newNode()
		node.children[r] = nextNode
		node = nextNode
	}
	node.isTerminal = true
	b.lastWord = word
	return nil
}

// minimize traverses up from the end of the last word's path, replacing
// nodes with equivalent ones in registeredNodes and registering new nodes.
func (b *DAWGBuilder) minimize(downTo int) {
	// Traverse path of last word from its end to the common prefix
	path := make([]*DAWGNode, len(b.lastWord)+1)
	path[0] = b.root
	for i, r := range b.lastWord {
		path[i+1] = path[i].children[r]
	}

	for i := len(b.lastWord); i > downTo; i-- {
		parent := path[i-1]
		child := path[i]
		char := rune(b.lastWord[i-1])

		sig := child.signature()
		if existingNode, ok := b.registeredNodes[sig]; ok {
			// Replaces child with node that already exists
			parent.children[char] = existingNode
		} else {
			// Register unique node otherwise
			b.registeredNodes[sig] = child
		}
	}
}

// Finish minimizes the last word added, and returns the immutable DAWG.
func (b *DAWGBuilder) Finish() *DAWG {
	b.minimize(0) // Minimize last word
	return &DAWG{root: b.root}
}
