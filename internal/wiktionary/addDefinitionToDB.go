package wiktionary

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/pbojar/dictextract/internal/database"
)

func addDefinitionToDB(word, pos, def string, dbQueries *database.Queries) (added bool, err error) {

	// Ensure word and pos are lowercase
	word = strings.ToLower(word)
	pos = strings.ToLower(pos)

	// Attempt to find existing entry in words
	wordID, err := dbQueries.GetIDByWord(context.Background(), word)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	// Add word if not found
	if wordID == 0 {
		dbWord, err := dbQueries.CreateWord(context.Background(), word)
		if err != nil {
			return false, err
		}
		wordID = dbWord.ID
	}

	// Attempt to find existing entry in pos
	posID, err := dbQueries.GetIDByPos(context.Background(), pos)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	// Add pos if not found
	if posID == 0 {
		dbPos, err := dbQueries.CreatePos(context.Background(), pos)
		if err != nil {
			return false, err
		}
		posID = dbPos.ID
	}

	// Return if definition already exists
	defExists, err := dbQueries.DefinitionExists(context.Background(), database.DefinitionExistsParams{
		WordID: wordID,
		PosID:  posID,
	})
	if err != nil {
		return false, err
	}
	if defExists {
		return false, nil
	}

	// Add definition
	_, err = dbQueries.CreateDefinition(context.Background(), database.CreateDefinitionParams{
		WordID:     wordID,
		PosID:      posID,
		Definition: def,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
