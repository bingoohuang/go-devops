package main

import (
	"encoding/json"
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

func WriteDbJson(db *badger.DB, k string, v interface{}, dur time.Duration) {
	j, _ := json.Marshal(v)
	WriteDb(db, k, j, dur)
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

func ReadDbJson(db *badger.DB, k string, v interface{}) (bool, error) {
	result, e := ReadDb(db, k)
	if e != nil {
		return false, e
	}

	if result == nil {
		return false, nil
	}

	e = json.Unmarshal(result, v)
	return true, e
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
