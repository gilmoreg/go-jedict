package main

import (
	"fmt"

	"github.com/gilmoreg/go-jedict/internal/ipa"
	"github.com/gilmoreg/go-jedict/internal/jedict"
)

func main() {
	ipaDict, err := ipa.New()
	if err != nil {
		panic(err)
	}
	jeDict, err := jedict.New()
	if err != nil {
		panic(err)
	}
	// // jEntry := jeDict.Entries[93587]
	withoutEntries := make([]jedict.Entry, 0)

	for _, e := range jeDict.Entries {
		if len(e.Readings) <= 0 && len(e.Kanji) <= 0 {
			fmt.Printf("entry %v has no readings or kanji\n", e.Sequence)
			continue
		}
		ipaEntries := make([]ipa.Entry, 0)
		if len(e.Kanji) > 0 {
			for _, k := range e.Kanji {
				ipaEntries = append(ipaEntries, ipaDict.FindEntriesBySurface(k)...)
			}
		}
		if len(ipaEntries) <= 0 {
			for _, r := range e.Readings {
				ipaEntries = append(ipaEntries, ipaDict.FindEntriesByReading(r)...)
			}
		}

		if len(ipaEntries) <= 0 {
			withoutEntries = append(withoutEntries, e)
		}
	}

	for i := 0; i < 20; i++ {
		fmt.Printf("Example: %v\n", withoutEntries[i])
	}
	percent := float32(len(withoutEntries)) / float32(len(jeDict.Entries))
	fmt.Printf("JEDict entries without IPA entries: %v, or %v percent\n", len(withoutEntries), percent)
}
