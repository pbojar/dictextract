package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pbojar/dictextract/internal/database"
	"github.com/pbojar/dictextract/internal/dawg"
	"github.com/pbojar/dictextract/internal/wiktionary"
)

type cliCommand struct {
	name        string
	description string
	callback    func(s *state, args ...string) error
}

func makeCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays help message.",
			callback:    commandHelp,
		},
		"lsRaws": {
			name:        "lsRaws",
			description: "Lists saved raw files.",
			callback:    commandListRaws,
		},
		"lsDAWGs": {
			name:        "lsDAWGs",
			description: "Lists saved DAWGs.",
			callback:    commandListDAWGs,
		},
		"makeDB": {
			name:        "makeDB <rawFileName>",
			description: "Makes a DB from words and definitions extracted from <rawFileName>.",
			callback:    commandMakeDB,
		},
		"makeDAWG": {
			name: "makeDAWG <minWordLen> <maxWordLen> <saveFileName>",
			description: `Makes a DAWG from words with lengths between <minWordLen> and <maxWordLen> (inclusive) 
    found in the current database. Saves the DAWG as a .gob file in the configured save directory.`,
			callback: commandMakeDAWG,
		},
	}
	return commands
}

func commandHelp(s *state, args ...string) error {
	fmt.Println(`dictextract is a command line tool to build databases of words and definitions and 
Directed Acyclic Word Graphs (DAWGs) from open source dictionaries (e.g., wiktionary) for use in word games.`)
	fmt.Print("\nUsage:\n")
	commands := makeCommands()
	for _, cliComm := range commands {
		fmt.Printf("  %s\n    %s\n", cliComm.name, cliComm.description)
	}
	return nil
}

func commandListDAWGs(s *state, args ...string) error {
	dawgDir := *s.cfg.DAWGSaveDirPath
	files, err := os.ReadDir(dawgDir)
	if err != nil {
		return err
	}
	if len(files) > 0 {
		fmt.Printf("The following DAWGs were found in '%s'...\n", dawgDir)
		for _, file := range files {
			fmt.Printf("  %s\n", file.Name())
		}
	} else {
		fmt.Printf("No DAWGs found in '%s'\n", dawgDir)
	}
	return nil
}

func commandListRaws(s *state, args ...string) error {
	rawDir := *s.cfg.RawDictDirPath
	files, err := os.ReadDir(rawDir)
	if err != nil {
		return err
	}
	if len(files) > 0 {
		fmt.Printf("The following files were found in '%s'...\n", rawDir)
		for _, file := range files {
			fmt.Printf("  %s\n", file.Name())
		}
	} else {
		fmt.Printf("No files found in '%s'\n", rawDir)
	}
	return nil
}

func commandMakeDB(s *state, args ...string) error {
	gzFilepath := args[0]
	err := wiktionary.ExtractToDB(gzFilepath, s.db)
	if err != nil {
		return err
	}
	return nil
}

func commandMakeDAWG(s *state, args ...string) error {

	// Check for proper number of args
	if len(args) != 3 {
		return fmt.Errorf("error: expected 3 arguments, '%d' given", len(args))
	}

	// Ensure first two args can be converted to integers
	minLenStr := args[0]
	maxLenStr := args[1]
	minLen, err := strconv.Atoi(minLenStr)
	if err != nil {
		return fmt.Errorf("error: '%s' is not convertable to an integer", minLenStr)
	}
	maxLen, err := strconv.Atoi(maxLenStr)
	if err != nil {
		return fmt.Errorf("error: '%s' is not convertable to an integer", maxLenStr)
	}

	// Ensure minLen < maxLen
	if minLen >= maxLen {
		return fmt.Errorf("error: <minLen> must be less than <maxLen>")
	}

	// Check for DAWGSaveDir and existing file name
	dawgDir := s.cfg.DAWGSaveDirPath
	if _, err := os.Stat(*dawgDir); os.IsNotExist(err) {
		return fmt.Errorf("error: directory '%s' does not exist", *dawgDir)
	}
	saveFileName := args[2] + ".gob"
	dawgSavePath := filepath.Join(*dawgDir, saveFileName)
	_, err = os.Stat(dawgSavePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return fmt.Errorf("error: file '%s' already exists", dawgSavePath)
	}

	// Get sorted words within range from DB
	fmt.Print("Getting words from db... ")
	sortedWords, err := s.db.GetWordsWithLenInRangeSorted(context.Background(), database.GetWordsWithLenInRangeSortedParams{
		Minlen: fmt.Sprintf("%d", minLen),
		Maxlen: fmt.Sprintf("%d", maxLen),
	})
	if err != nil {
		return fmt.Errorf("error: could not get words from db\n%v", err)
	}
	totalWords := len(sortedWords)
	fmt.Printf("Done!\nFound %d words!\n\n", totalWords)

	// Make DAWG from sorted words list
	fmt.Println("Building DAWG...")
	builder := dawg.NewDAWGBuilder()
	for i, w := range sortedWords {
		err := builder.Insert(w)
		if err != nil {
			return fmt.Errorf("error: could not insert '%s'", w)
		}
		fmt.Printf("\033[2K\r%d of %d words added to DAWG", i+1, totalWords)
	}
	finalDAWG := builder.Finish()
	fmt.Printf("\nDone!\n\n")

	// Save DAWG to file
	fmt.Printf("Saving DAWG to '%s'... ", dawgSavePath)
	err = finalDAWG.SaveAsGob(dawgSavePath)
	if err != nil {
		return err
	}
	fmt.Printf("Done!\n")

	return nil
}
