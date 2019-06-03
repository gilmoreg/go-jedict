package jedict

// Dictionary -
type Dictionary struct {
	Entries []Entry
}

// Entry -
type Entry struct {
	// Unique identifier per record
	Sequence int

	Kanji []string

	// Readings in kana
	Readings []string

	// List of meanings in English
	Meanings []Meaning
}

// Meaning -
type Meaning struct {
	PartOfSpeech []string
	Gloss        string
	Misc         []string
	IPA          string
}

// New - create new JEDict dictionary
func New() (*Dictionary, error) {
	entries, err := readXML("data/JMdict_e.xml")
	if err != nil {
		return nil, err
	}
	return &Dictionary{Entries: entries}, nil
}
