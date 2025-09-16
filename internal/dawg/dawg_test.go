package dawg

import (
	"sort"
	"testing"
)

func TestDAWGNodeSignature(t *testing.T) {

	// Create a DAWGNode
	testWord := "doggy"

	// Start with root
	root := DAWGNode{
		id:       0,
		children: make(map[rune]*DAWGNode),
	}

	// Add children
	node := &root
	for i, r := range testWord {
		nextNode := &DAWGNode{
			id:       i + 1,
			children: make(map[rune]*DAWGNode),
		}
		node.children[r] = nextNode
		node = nextNode
	}
	node.isTerminal = true

	// Return to root and get node signatures
	node = &root
	sigs := make([]string, len(testWord)+1)
	for _, r := range testWord {
		sigs[node.id] = node.signature()
		node = node.children[r]
		if node.isTerminal {
			sigs[node.id] = node.signature()
		}
	}

	// Compare to expected signatures
	expectedSigs := []string{"0_d1_", "0_o2_", "0_g3_", "0_g4_", "0_y5_", "1_"}
	for i, expectedSig := range expectedSigs {
		if sigs[i] != expectedSig {
			t.Errorf("Signature of node %d does not match: got '%s', want '%s'", i, sigs[i], expectedSig)
		}
	}
}

func TestDAWGBuilderInsert(t *testing.T) {
	unsortedWords := []string{"zoo", "moo"}
	sortedWords := []string{
		"cat",
		"cats",
		"catch",
	}
	sort.Strings(sortedWords)

	tests := []struct {
		name    string
		words   []string
		builder *DAWGBuilder
		wantErr bool
	}{
		{
			name:    "Unsorted",
			words:   unsortedWords,
			builder: NewDAWGBuilder(),
			wantErr: true,
		},
		{
			name:    "Sorted",
			words:   sortedWords,
			builder: NewDAWGBuilder(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var err error
			for _, w := range tt.words {
				err = tt.builder.Insert(w)
				if err != nil {
					break
				}
			}

			// Catch unwanted errors from Insert loop
			if (err != nil) != tt.wantErr {
				t.Errorf("DAWGBuilder.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDAWGBuilderFinalize(t *testing.T) {
	fewWords := []string{
		"cat",
		"cats",
		"catch",
	}
	sort.Strings(fewWords)
	manyWords := []string{
		"cat",
		"car",
		"cats",
		"catch",
		"cache",
		"dog",
		"dogs",
		"doggy",
	}
	sort.Strings(manyWords)

	tests := []struct {
		name         string
		words        []string
		builder      *DAWGBuilder
		expectedKeys []string
	}{
		{
			name:    "Few words",
			words:   fewWords,
			builder: NewDAWGBuilder(),
			expectedKeys: []string{
				"0_a2_",
				"0_t3_",
				"0_h5_",
				"1_c4_s5_",
				"1_",
			},
		},
		{
			name:    "Many words",
			words:   manyWords,
			builder: NewDAWGBuilder(),
			expectedKeys: []string{
				"0_a2_",
				"0_h4_",
				"0_h5_",
				"0_e5_",
				"0_y5_",
				"0_c3_r5_t7_",
				"1_c8_s5_",
				"0_o12_",
				"0_g13_",
				"1_g14_s5_",
				"1_",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range tt.words {
				err := tt.builder.Insert(w)
				// Catch unwanted errors from Insert loop
				if err != nil {
					t.Errorf("DAWGBuilder.Insert() error = %v", err)
					return
				}
			}

			// Finalize DAWG
			tt.builder.Finish()

			if len(tt.builder.registeredNodes) != len(tt.expectedKeys) {
				t.Errorf("Mismatch in size of registeredNodes")
			}
			for _, expectedKey := range tt.expectedKeys {
				_, ok := tt.builder.registeredNodes[expectedKey]
				if !ok {
					t.Errorf("registeredNodes is missing expected key '%s'", expectedKey)
				}
			}
		})
	}
}

func TestDAWGContains(t *testing.T) {
	testWords := []string{
		"cat",
		"car",
		"cats",
		"catch",
		"cache",
		"dog",
		"dogs",
		"doggy",
	}
	sort.Strings(testWords)
	builder := NewDAWGBuilder()
	for _, w := range testWords {
		err := builder.Insert(w)
		// Catch unwanted errors from Insert loop
		if err != nil {
			t.Errorf("DAWGBuilder.Insert() error = %v", err)
			return
		}
	}
	dawg := builder.Finish()

	tests := []struct {
		name     string
		word     string
		dawg     *DAWG
		expected bool
	}{
		{
			name:     "DAWG contains dog",
			word:     "dog",
			dawg:     dawg,
			expected: true,
		},
		{
			name:     "DAWG doesn't contain do",
			word:     "do",
			dawg:     dawg,
			expected: false,
		},
		{
			name:     "DAWG doesn't contain doggo",
			word:     "doggo",
			dawg:     dawg,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dawg.Contains(tt.word)
			if result != tt.expected {
				t.Errorf("dawg.Contains(%s) returned %t, expected %t", tt.word, result, tt.expected)
			}
		})
	}
}

func TestDAWGStartsWith(t *testing.T) {
	testWords := []string{
		"cat",
		"car",
		"cats",
		"catch",
		"cache",
		"dog",
		"dogs",
		"doggy",
	}
	sort.Strings(testWords)
	builder := NewDAWGBuilder()
	for _, w := range testWords {
		err := builder.Insert(w)
		// Catch unwanted errors from Insert loop
		if err != nil {
			t.Errorf("DAWGBuilder.Insert() error = %v", err)
			return
		}
	}
	dawg := builder.Finish()

	tests := []struct {
		name     string
		prefix   string
		dawg     *DAWG
		expected bool
	}{
		{
			name:     "DAWG words start with dog",
			prefix:   "dog",
			dawg:     dawg,
			expected: true,
		},
		{
			name:     "DAWG words start with car",
			prefix:   "car",
			dawg:     dawg,
			expected: true,
		},
		{
			name:     "DAWG words don't start with cars",
			prefix:   "cars",
			dawg:     dawg,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dawg.StartsWith(tt.prefix)
			if result != tt.expected {
				t.Errorf("dawg.StartsWith(%s) returned %t, expected %t", tt.prefix, result, tt.expected)
			}
		})
	}
}
