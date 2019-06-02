package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

// NewBinaryStorageProvider - new binary storage provider
func NewBinaryStorageProvider() *Binary {
	return &Binary{
		entries: make([]Entry, 0),
	}
}

// Binary - write entries to binary file
type Binary struct {
	entries []Entry
}

// StoreEntry - add entry to commit batch
func (b *Binary) StoreEntry(e Entry) error {
	b.entries = append(b.entries, e)
	return nil
}

// UncommittedEntries - number of uncommitted entries
func (b *Binary) UncommittedEntries() int {
	return len(b.entries)
}

// Commit - write to file
func (b *Binary) Commit() error {
	// size := binary.Size(b.entries)
	// buf := make([]byte, size)
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(b.entries)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of entries written: %v; size: %v", len(b.entries), binary.Size(b.entries))
	err = ioutil.WriteFile("binarytest.bin", buffer.Bytes(), 0644)
	return err
}

// Purge - no op to fulfill interface
func (b *Binary) Purge() error {
	return nil
}
