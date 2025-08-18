package main

import (
	"log"
)

func main() {
	// wiktionENURL := "https://kaikki.org/dictionary/raw-wiktextract-data.jsonl.gz"

	gzFilepath := "/tmp/dict-949316240"

	// gzFilepath, err := downloadDictionary(wiktionENURL)
	// if err != nil {
	// 	log.Fatalf("Error downloading dictionary: %v", err)
	// }
	// fmt.Println(gzFilepath)

	err := extractValidWords(gzFilepath)
	if err != nil {
		log.Fatalf("Error extracting and building EN dictionary: %v", err)
	}

}
