package wiktionary

import (
	"slices"
	"strings"
	"unicode"
)

// Unicode range table for A-Z and a-z
var engAlphaRange = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 'A', Hi: 'Z', Stride: 1},
		{Lo: 'a', Hi: 'z', Stride: 1},
	},
}

// Unicode range table for A-Z
var engCAPSRange = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 'A', Hi: 'Z', Stride: 1},
	},
}

// isEngAlphaOnly returns true if all runes in the string are in engAlphaRange (A-Z and a-z).
// unicode.RangeTable is used here in hopes of supporting runes with accents for other languages.
func isEnAlphaOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsOneOf([]*unicode.RangeTable{engAlphaRange}, r) {
			return false
		}
	}
	return true
}

// hasInitialism returns true if the length of a sequence of capital letters in
// 's' meets or exceeds the 'cutOff' and false otherwise.
func hasInitialism(s string, cutOff int) bool {
	count := 0
	for _, r := range s {
		if unicode.IsOneOf([]*unicode.RangeTable{engCAPSRange}, r) {
			count += 1
		} else {
			count = 0
		}
		if count >= cutOff {
			return true
		}
	}
	return false
}

// Filter returns true if a wiktionary entry, represented by the struct wiktionLite,
// adheres to the following requirements:
//  1. The language code is "en".
//  2. The part of speech is noun, pron(oun), verb, adj(ective),
//     adv(erb), prep(osition), conj(unction), or int(er)j(ection).
//  3. At least one definition exists.
//  4. The word contains only english alphabet runes.
//  5. The word is not an initialism (per hasInitialism)
//
// and returns false otherwise.
func filter(w *wiktionLite) bool {

	// Match lang code
	if w.LangCode != "en" {
		return false
	}

	// Ensure word is an accepted part of speech
	acceptedPos := []string{"noun", "pron", "verb", "adj", "adv", "prep", "conj", "intj"}
	if !slices.Contains(acceptedPos, w.Pos) {
		return false
	}

	// Check for definition, invalidate words without definitions
	if len(w.Senses) == 0 || len(w.Senses[0].Glosses) == 0 {
		return false
	}

	// Check that word contains only English Alphabet Letters
	if !isEnAlphaOnly(w.Word) {
		return false
	}

	// Check if word is an initialism/acronym by checking if it has adjacent CAPS
	if hasInitialism(w.Word, 2) {
		return false
	}

	// Finally, checks if word is an initialism/acronym by checking the definition for initialism/acronym
	lowerDef := strings.ToLower(w.Senses[0].Glosses[0])
	if strings.Contains(lowerDef, "initialism") || strings.Contains(lowerDef, "acronym") {
		return false
	}

	return true
}
