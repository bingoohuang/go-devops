package main

import (
	"github.com/dgraph-io/badger"
	"log"
	"time"
)

func OpenDb(dataDir string) *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = dataDir
	opts.ValueDir = dataDir
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
	//defer db.Close()
}

func WriteDb(db *badger.DB, k string, v []byte, dur time.Duration) {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.SetWithTTL([]byte(k), v, dur)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func ReadDb(db *badger.DB, k string) (result []byte, err error) {
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		result = val
		return nil
	})
	return
}
