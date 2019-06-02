package main

/*
IPA Notes

最初の4カラム目までは, 必須項目で, (first 4 columns are required)

表層形 (単語そのもの) (surface)
左連接状態番号 (left link state number)
右連接状態番号 (right link state number)
コスト (cost)

5カラム目以降は「素性」と呼ばれる項目です. (5th and beyond are called "features")
5カラム目は品詞, 6カラム目は品詞再分 類等) (5th column part of speech, 6th pos classification)

(例: 品詞, 品詞細分類, 活用型, 活用形, 原形, 読み, 発音)
part of speech, subdivision of parts of speech, inflection type, inflection type, original form, reading, pronunciation

So we have:
0 surface
1 link left state number
2 link right state number
3 cost
4-9 part of speech
10 base
11 reading
12 pronunciation

*/

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// IPAEntry from CSV file
// ex よ,1,1,6514,その他,間投,*,*,*,*,よ,ヨ,ヨ
type IPAEntry struct {
	Surface string   `json:"surface"`
	POS     []string `json:"pos"`
	Base    string   `json:"base"`
	Reading string   `json:"reading"`
	Pron    string   `json:"pron"`
}

var csvFile = regexp.MustCompile(`.+\.csv$`)

func main() {
	entries, err := LoadIPADictionary()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Number of entries read: %v, example: %v\n", len(entries), entries[200000])
}

// LoadIPADictionary -
func LoadIPADictionary() ([]IPAEntry, error) {
	var files []string
	var ipaEntries = make([]IPAEntry, 0)

	err := filepath.Walk("data/ipa", func(path string, info os.FileInfo, err error) error {
		if csvFile.MatchString(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		entries, err := loadIPACSV(file)
		if err != nil {
			return nil, err
		}
		ipaEntries = append(ipaEntries, entries...)
	}

	return ipaEntries, nil
}

func loadIPACSV(path string) ([]IPAEntry, error) {
	fmt.Printf("Reading %s\n", path)
	csvFile, err := os.Open(path)
	defer csvFile.Close()
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var entries []IPAEntry
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		entries = append(entries, IPAEntry{
			Surface: line[0],
			POS:     line[4:9],
			Base:    line[10],
			Reading: line[11],
			Pron:    line[12],
		})
	}
	return entries, nil
}
