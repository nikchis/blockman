// Copyright (c) 2018 Nikita Chisnikov
// Distributed under the MIT/X11 software license

package main

import (
	"log"
	"time"
	"errors"
	"github.com/boltdb/bolt"
)

func setData(dbFile, data string) error {
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println("db.Open():", err.Error())
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("block"))
		if err != nil {
			return err
		}
		if err := b.Put([]byte("data"), []byte(data)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("db.Update():", err.Error())
		return err
	}
	return nil
}

func getData(dbFile string) (sres string, err error) {
	if (dbFile == "") { return sres, errors.New("empty path") }
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println("db.Open():", err.Error())
		return sres, err
	}
	defer db.Close()
	var res []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("block"))
		if b == nil { 
			return errors.New("bucket is nil")
		}
		res = b.Get([]byte("data"))
		if res == nil {
			return errors.New("value is nil")
		}
		return nil
	})
	if err != nil {
		log.Println("db.View():", err.Error())
		return sres, err
	}
	return string(res), nil
}
