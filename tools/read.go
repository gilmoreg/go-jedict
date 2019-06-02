package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"

	"github.com/gilmoreg/go-jedict/dictionary/storage"
)

func main() {
	binbytes, err := ioutil.ReadFile("data/binarytest.bin")
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(binbytes)
	dec := gob.NewDecoder(buf)
	var entries []storage.Entry
	err = dec.Decode(&entries)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of entries read: %v, example: %v", len(entries), entries[1])
}
