package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// MongoDBStorageProvider - MongoDBStorageProvider
type MongoDBStorageProvider struct {
	connectionString string
	client           *mongo.Client
	collection       *mongo.Collection
}

// NewMongoDBStorageProvider - create a new MongoDBStorageProvider
func NewMongoDBStorageProvider(connectionString string) *MongoDBStorageProvider {
	sp := new(MongoDBStorageProvider)
	sp.connectionString = connectionString
	client, err := mongo.Connect(context.TODO(), connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	sp.client = client
	sp.collection = client.Database("jedict").Collection("entries")
	return sp
}

// LookupKanji - Look up a word by kanji expression
func (m MongoDBStorageProvider) LookupKanji(string) (Entry, error) {
	return Entry{}, errors.New("Not yet implemented")
}

// LookupReading - Look up a word by hiragana/katakana expression
func (m MongoDBStorageProvider) LookupReading(string) (Entry, error) {
	return Entry{}, errors.New("Not yet implemented")
}

// LookupMeaning - Look up a word by English expression
func (m MongoDBStorageProvider) LookupMeaning(string) (Entry, error) {
	return Entry{}, errors.New("Not yet implemented")
}

// LookupWord - Look up a word by any of kanji, reading, or English expression
// (returns the first result in that order)
func (m MongoDBStorageProvider) LookupWord(string) (Entry, error) {
	return Entry{}, errors.New("Not yet implemented")
}

// StoreEntry - Store an entry in the database.
func (m MongoDBStorageProvider) StoreEntry(e Entry) error {
	fmt.Println("Storing entry", e.Sequence)
	_, err := m.collection.InsertOne(context.TODO(), e)
	return err
}

// UncommittedEntries - Returns the number of entries that have not been committed to disk
func (m MongoDBStorageProvider) UncommittedEntries() int {
	return 0
}

// Commit - Save all changes to the database
func (m MongoDBStorageProvider) Commit() error {
	return nil
}

// Purge all entries from this storage (DESTRUCTIVE!)
func (m MongoDBStorageProvider) Purge() error {
	return errors.New("Unimplemented")
}
