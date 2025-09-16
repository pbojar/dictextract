package dawg

// DAWG is an immutable Directed Acyclic Word Graph.
type DAWG struct {
	root *DAWGNode
}

// Contains checks if a word exists in the DAWG.
func (d *DAWG) Contains(word string) bool {
	node := d.root
	for _, r := range word {
		child, ok := node.children[r]
		if !ok {
			return false
		}
		node = child
	}
	return node.isTerminal
}

// StartsWith checks if any word in the DAWG starts with the given prefix.
func (d *DAWG) StartsWith(prefix string) bool {
	node := d.root
	for _, r := range prefix {
		child, ok := node.children[r]
		if !ok {
			return false
		}
		node = child
	}
	return true
}
