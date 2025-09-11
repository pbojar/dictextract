package wiktionary

type wiktionLite struct {
	Senses []struct {
		Glosses []string `json:"glosses"`
	} `json:"senses"`
	Pos      string `json:"pos"`
	Word     string `json:"word"`
	LangCode string `json:"lang_code"`
}
