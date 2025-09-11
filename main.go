package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pbojar/dictextract/internal/database"
	"github.com/pbojar/dictextract/internal/wiktionary"
)

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	gzFilepath := os.Getenv("RAW_DICT_PATH")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	err = wiktionary.ExtractToDB(gzFilepath, dbQueries)
	if err != nil {
		log.Fatalf("Error extracting and building EN dictionary: %v", err)
	}

}
