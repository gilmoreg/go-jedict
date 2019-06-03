package ipa

import (
	"strings"
)

// Dictionary - IPA Dictionary
type Dictionary struct {
	entriesBySurface map[string][]Entry
	entriesByReading map[string][]Entry
}

// New - return new IPA Dictionary
func New() (*Dictionary, error) {
	entriesBySurface, entriesByReading, err := loadIPADictionary()
	if err != nil {
		return nil, err
	}
	return &Dictionary{
		entriesBySurface: entriesBySurface,
		entriesByReading: entriesByReading,
	}, nil
}

// FindEntriesBySurface - find matching entries by surface
func (d *Dictionary) FindEntriesBySurface(surface string) []Entry {
	if val, ok := d.entriesBySurface[surface]; ok {
		return val
	}
	return make([]Entry, 0)
}

// FindEntriesByReading - find matching entries by reading
func (d *Dictionary) FindEntriesByReading(reading string) []Entry {
	if val, ok := d.entriesByReading[toKatakana(reading)]; ok {
		return val
	}
	return make([]Entry, 0)
}

// https://www.lemoda.net/go/hiragana-to-katakana/index.html
func toKatakana(hira string) string {
	result := strings.Map(hira2kata, hira)
	return result
}

func hira2kata(hira rune) rune {
	if (hira >= 'ぁ' && hira <= 'ゖ') || (hira >= 'ゝ' && hira <= 'ゞ') {
		return hira + 0x60
	}
	return hira
}
