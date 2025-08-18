package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type wiktionLite struct {
	Senses []struct {
		Glosses []string `json:"glosses"`
	} `json:"senses"`
	Pos      string `json:"pos"`
	Word     string `json:"word"`
	LangCode string `json:"lang_code"`
}

type definition struct {
	Def string `json:"definition"`
	Pos string `json:"pos"`
}

type dictEntry struct {
	Word        string       `json:"word"`
	Definitions []definition `json:"definitions"`
}

var engAlphaRange = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 'A', Hi: 'Z', Stride: 1},
		{Lo: 'a', Hi: 'z', Stride: 1},
	},
}

var engCAPSRange = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 'A', Hi: 'Z', Stride: 1},
	},
}

func isEngAlphaOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsOneOf([]*unicode.RangeTable{engAlphaRange}, r) {
			return false
		}
	}
	return true
}

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

func (w *wiktionLite) isValid(langCode string, minLength int, maxLength int) bool {

	// Match lang code
	if w.LangCode != langCode {
		return false
	}

	// Exclude names
	if w.Pos == "name" {
		return false
	}

	// Check for definition, invalidate words without definitions
	if len(w.Senses) == 0 || len(w.Senses[0].Glosses) == 0 {
		return false
	}

	// Check that word contains only English Alphabet Letters and is within min/maxLength (inclusive)
	if !isEngAlphaOnly(w.Word) || len(w.Word) < minLength || len(w.Word) > maxLength {
		return false
	}

	// Check if word is an initialism/acronym by checking if it has adjacent CAPS
	isInitialism := hasInitialism(w.Word, 2)

	return !isInitialism
}

func extractValidWords(gzFilepath string) (err error) {
	langCode := "en"

	// Open compressed file for reading
	tmpFile, err := os.Open(gzFilepath)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// Create GZIP reader
	gzReader, err := gzip.NewReader(tmpFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// Create langCode directory
	err = os.Mkdir(langCode, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	// Create dict map
	dictMap := make(map[string][]definition)

	// Use buffered scanner to read line by line
	scanner := bufio.NewScanner(gzReader)
	const maxCapacity int = 1 << 24 // Lines in the EN wiktionary are very long
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	const minWordLen int = 3
	const maxWordLen int = 18
	for scanner.Scan() {

		// Unmarshal json
		var entry wiktionLite
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Println("Error parsing JSON: ", err)
			continue
		}

		// Ensure word does not contain invalid characters, matches langCode, is longer than 3 letters,
		// and is not an initialism
		if !entry.isValid(langCode, minWordLen, maxWordLen) {
			continue
		}

		key := strings.ToLower(entry.Word)
		dictMap[key] = append(dictMap[key], definition{
			Def: entry.Senses[0].Glosses[0],
			Pos: entry.Pos,
		})

	}

	if err = scanner.Err(); err != nil {
		return err
	}

	// Create index output file
	indexFile, err := os.Create(langCode + "/index.txt")
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// Create dict output file
	dictFile, err := os.Create(langCode + "/dict.jsonl")
	if err != nil {
		return err
	}
	defer dictFile.Close()

	// Write to dict.jsonl and index.txt
	for word, definitions := range dictMap {

		// Output full entry to langCode/dict.jsonl
		dat, err := json.Marshal(dictEntry{
			Word:        word,
			Definitions: definitions,
		})
		if err != nil {
			return err
		}
		_, err = dictFile.Write(dat)
		if err != nil {
			return err
		}
		_, err = dictFile.WriteString("\n")
		if err != nil {
			return err
		}

		// Output word to langCode/index.txt
		_, err = indexFile.WriteString(word + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
