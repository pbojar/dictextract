package wiktionary

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pbojar/dictextract/internal/database"
)

// ExtractToDB decompresses gzFilepath line-by-line and attempts to parse each line as a json following the wiktionLite
// structure. Entries are then filtered and the first definition for the word, pos pair is added to the database.
func ExtractToDB(gzFilepath string, dbQueries *database.Queries) (err error) {

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

	// Use buffered scanner to read line by line
	numAdded := 0
	numFiltered := 0
	numDupes := 0
	scanner := bufio.NewScanner(gzReader)
	const maxCapacity int = 1 << 24 // Lines in the EN wiktionary are very long
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	fmt.Println("Extracting and Adding Definitions...")
	for scanner.Scan() {

		// Unmarshal json
		var entry wiktionLite
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			log.Printf("Error parsing JSON: %s", err)
			continue
		}

		if !filter(&entry) {
			numFiltered++
			continue
		}

		added, err := addDefinitionToDB(
			entry.Word, entry.Pos, entry.Senses[0].Glosses[0],
			dbQueries,
		)
		if err != nil {
			return err
		}
		if added {
			numAdded++
		} else {
			numDupes++
		}
		fmt.Printf("\033[2K\rDefinitions (added, filtered, dupes): (%d, %d, %d)", numAdded, numFiltered, numDupes)

	}
	fmt.Printf("\nExtract and add complete!\n")

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
