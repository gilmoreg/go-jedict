package jedict

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type jmDictXMLDoc struct {
	JmDict jmDict `xml:"JMdict"`
}

type jmDict struct {
	XMLName xml.Name `xml:"JMdict"`
	Entries []entry  `xml:"entry"`
}

type entry struct {
	XMLName xml.Name `xml:"entry"`
	EntSeq  int      `xml:"ent_seq"`
	KEle    []kEle   `xml:"k_ele"`
	REle    []rEle   `xml:"r_ele"`
	Sense   []sense  `xml:"sense"`
}

type entity struct {
	XMLName xml.Name `xml:"!ENTITY"`
	content string   `xml:,innerxml`
}

type rEle struct {
	Reb string `xml:"reb"`
}

type kEle struct {
	Keb string `xml:"keb"`
}

type sense struct {
	Gloss []Gloss  `xml:"gloss"`
	Pos   []string `xml:"pos"`
	Misc  []string `xml:"misc"`
}

// Gloss -
type Gloss struct {
	Content string `xml:",chardata"`
	Lang    string `xml:"lang,attr"`
}

// Converts an XML read entity into an Entry
func (e *entry) convertToEntry(edictMapping map[string]string) Entry {
	entry := Entry{
		Sequence: e.EntSeq,
		Kanji:    make([]string, 0),
		Readings: make([]string, 0),
		Meanings: make([]Meaning, 0),
	}
	for _, reading := range e.REle {
		entry.Readings = append(entry.Readings, reading.Reb)
	}
	for _, kanji := range e.KEle {
		entry.Kanji = append(entry.Kanji, kanji.Keb)
	}

	var currentPos string

	for _, sense := range e.Sense {
		pos := make([]string, 0)

		if len(sense.Pos) > 0 {
			for _, p := range sense.Pos {
				p = strings.ToLower(p)
				if val, ok := edictMapping[p]; ok {
					p = val
				}
				pos = append(pos, p)
				currentPos = p
			}
		} else {
			pos = append(pos, currentPos)
		}

		meanings := make([]string, 0)
		for _, gloss := range sense.Gloss {
			// For now, skip languages other than the default (english)
			if gloss.Lang == "" {
				meanings = append(meanings, gloss.Content)
			}
		}
		meaningsStr := strings.Join(meanings, ", ")
		misc := make([]string, 0)
		for _, m := range sense.Misc {
			misc = append(misc, m)
		}

		if meaningsStr != "" {
			entry.Meanings = append(entry.Meanings, Meaning{
				PartOfSpeech: pos,
				Gloss:        meaningsStr,
				Misc:         misc,
			})
		}
	}

	return entry
}

func readXML(filename string) ([]Entry, error) {
	fmt.Print("Loading JMDict...")
	reader, err := os.Open(filename)
	defer reader.Close()

	if err != nil {
		fmt.Printf("Error opening dictionary xml file: %s\n", err)
		return nil, err
	}

	edictMapFile, err := os.Open("data/edictmap.json")
	defer edictMapFile.Close()

	var edictMapping map[string]string
	edictRaw, _ := ioutil.ReadAll(edictMapFile)
	json.Unmarshal(edictRaw, &edictMapping)

	decoder := xml.NewDecoder(reader)
	decoder.Strict = false

	entries := make([]Entry, 0)

	// Load directives
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.Directive: // Reading an XML directive
			directive := token.(xml.Directive)
			results, err := findEntities(&directive)
			if err != nil {
				fmt.Printf("Entities could not be found: %s\n", err)
				continue
			}
			decoder.Entity = results
		case xml.StartElement: // Reading an XML element
			var entry entry
			if startElement.Name.Local == "entry" {
				decoder.DecodeElement(&entry, &startElement)
			}
			storeEntry := entry.convertToEntry(edictMapping)
			entries = append(entries, storeEntry)
			if err != nil {
				_ = fmt.Errorf("Error storing entry: %s\n %+v", err, entry)
				return nil, err
			}
		}
	}
	fmt.Println("complete.")
	return entries, nil
}

func findEntities(d *xml.Directive) (map[string]string, error) {
	directiveStr := string(*d)

	// Go's xml package doesn't parse the DOCTYPE directive, we
	// need to parse this manually. Since we're only looking for
	// ENTITY definitions, a regex will do.
	re, err := regexp.Compile("<!ENTITY\\s([0-9A-z]+)\\s\"(.+)\">")
	if err != nil {
		fmt.Printf("Error compiling regex: %s", err)
		return nil, err
	}

	entities := make(map[string]string)

	results := re.FindAllStringSubmatch(directiveStr, -1)
	for _, entity := range results {
		entities[entity[1]] = entity[2]
	}

	return entities, nil
}

func (dic *jmDict) ReadInto(reader io.Reader) error {
	decoder := xml.NewDecoder(reader)
	decoder.Strict = false

	// Load directives
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.Directive: // Reading an XML directive
			directive := token.(xml.Directive)
			results, err := findEntities(&directive)
			if err != nil {
				fmt.Printf("Entities could not be found: %s\n", err)
				continue
			}
			decoder.Entity = results
		case xml.StartElement: // Reading an XML element
			var entry entry
			if startElement.Name.Local == "entry" {
				decoder.DecodeElement(&entry, &startElement)
			}
			dic.Entries = append(dic.Entries, entry)
		}
	}
	return nil
}
