/*
`jedict` is the command line tool for importing the JMDICT XML dictionary
and performing lookups on the database.

Usage:
	jedict -db CONNECTION_STRING [-import XML_PATH]
		                         [-kanji EXPRESSION]
	                             [-meaning EXPRESSION]
	                             [-reading EXPRESSION]

XML_PATH is the location of the JMDICT XML file. See README.md for more details.

CONNECTION_STRING is the postgresql connection string for the database containing
the dictionary. The database must be prepopulated.

EXPRESSION is a kanji expression, english expression, or hiragana/katakana
expression respectively.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/gilmoreg/go-jedict/dictionary"
	"github.com/gilmoreg/go-jedict/dictionary/storage"
)

func doImport(p storage.Writer, xmlfile string) error {
	progress := make(chan float32)

	// Print out the progress while the import is running
	go func(progress chan float32) {
		incomplete := true
		var completion float32
		for incomplete == true {
			completion, incomplete = <-progress
			fmt.Printf("\rImport progress: %.2f%%", completion*100)
		}
		fmt.Printf("\n")
	}(progress)

	err := dictionary.ReadXMLIntoStorage(xmlfile, p, progress)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Connect to the dictionary database using the provided connetion string
	var xmlfile string
	var connectionString string
	flag.StringVar(&connectionString, "db", "", "MongoDB connection string")
	flag.StringVar(&xmlfile, "import", "", "JMdict file to import")
	flag.Parse()

	provider := storage.NewMongoDBStorageProvider(connectionString)

	// -import option for performing database import
	if xmlfile != "" {
		err := doImport(provider, xmlfile)
		if err != nil {
			_ = fmt.Errorf("error importing dictionary: %s", err)
			return
		}
	}
}
