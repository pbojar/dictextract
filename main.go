package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pbojar/dictextract/internal/config"
	"github.com/pbojar/dictextract/internal/database"
)

func main() {

	// Check at least one arg is given
	if len(os.Args) < 2 {
		fmt.Printf("error parsing args:\nusage: ./dictextract <command> (args...)\n")
		os.Exit(1)
	}

	// Get user command name and args
	user_cmd_name := os.Args[1]
	user_cmd_args := os.Args[2:]

	// Read user config
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
	}

	// Connect to current DB
	// TODO: Support changing DBs
	db, err := sql.Open("postgres", *cfg.DBURL)
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
	}
	dbQueries := database.New(db)

	// Initialize app state
	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Make commands
	c := makeCommands()

	// Check if user command is valid
	cliComm, exists := c[user_cmd_name]
	if !exists {
		// Exit if command does not exist
		log.Fatalf("Unknown command '%s'", user_cmd_name)
	} else {
		// Execute command callback function with user args
		err := cliComm.callback(&s, user_cmd_args...)
		if err != nil {
			log.Fatal(err)
		}
	}
}
