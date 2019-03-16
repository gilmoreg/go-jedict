package dictionary

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/gilmoreg/go-jedict/dictionary/storage"
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
	Gloss []Gloss `xml:"gloss"`
	Pos   string  `xml:"pos"`
}

// Gloss -
type Gloss struct {
	Content string `xml:",chardata"`
	Lang    string `xml:"lang,attr"`
}

// Converts an XML read entity into an entry ready for storage
func (e *entry) convertToStorageEntry() storage.Entry {
	sEntry := storage.Entry{
		Sequence: e.EntSeq,
		Kanji:    make([]string, 0),
		Readings: make([]string, 0),
		Meanings: make([]storage.Meaning, 0),
	}
	for _, reading := range e.REle {
		sEntry.Readings = append(sEntry.Readings, reading.Reb)
	}
	for _, kanji := range e.KEle {
		sEntry.Kanji = append(sEntry.Kanji, kanji.Keb)
	}
	for _, sense := range e.Sense {
		meanings := make([]string, 0)
		for _, gloss := range sense.Gloss {
			// For now, skip languages other than the default (english)
			if gloss.Lang == "" {
				meanings = append(meanings, gloss.Content)
			}
		}
		meaningsStr := strings.Join(meanings, ", ")
		if sense.Pos != "" && meaningsStr != "" {
			sEntry.Meanings = append(sEntry.Meanings, storage.Meaning{
				PartOfSpeech: sense.Pos,
				Gloss:        meaningsStr,
			})
		}
	}

	return sEntry
}

// ReadXMLIntoStorage -
func ReadXMLIntoStorage(filename string, provider storage.Writer, progress chan float32) error {
	reader, err := os.Open(filename)
	defer reader.Close()

	if err != nil {
		fmt.Printf("Error opening dictionary xml file: %s\n", err)
		return err
	}

	var fileSize int64
	if progress != nil {
		fileInfo, err := reader.Stat()
		if err != nil {
			return err
		}
		fileSize = fileInfo.Size()
		defer close(progress)
	}

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
			results, err := FindEntities(&directive)
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
			err = provider.StoreEntry(entry.convertToStorageEntry())
			if err != nil {
				fmt.Errorf("Error storing entry: %s\n %+v", err, entry)
				return err
			}
			if provider.UncommittedEntries() >= 1000 {
				if progress != nil {
					currentLocation := decoder.InputOffset()
					progress <- (float32(currentLocation) / float32(fileSize))
				}
				err = provider.Commit()
				if err != nil {
					fmt.Errorf("Error committing entries: %s", err)
					return err
				}
			}
		}
	}

	provider.Commit()
	if err != nil {
		fmt.Errorf("Error committing entries: %s", err)
		return err
	}

	if progress != nil {
		// 100% complete
		progress <- 1.0
	}

	return nil
}

// FindEntities -
func FindEntities(d *xml.Directive) (map[string]string, error) {
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
			results, err := FindEntities(&directive)
			if err != nil {
				fmt.Printf("Entities could not be found: %s\n", err)
				continue
			}
			//fmt.Printf("Entities: %+v\n", results)
			decoder.Entity = results
		case xml.StartElement: // Reading an XML element
			var entry entry
			if startElement.Name.Local == "entry" {
				decoder.DecodeElement(&entry, &startElement)
			}
			dic.Entries = append(dic.Entries, entry)
		}
	}

	//return decoder.Decode(&dict)
	return nil

}
