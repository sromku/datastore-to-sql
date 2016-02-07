package main

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"io"
	"io/ioutil"
	"log"
	"os"
	datastore "./datastore"
	ds "./internal/datastore"
)

var Strict bool = true

// Load the backup into real model
func LoadBackup(backupFilePath string, dst interface{}, onPreload func(dst interface{}), onResult func(dst interface{})) {
	f, err := os.Open(backupFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	journals := journal.NewReader(f, nil, Strict, true)
	for {
		j, err := journals.Next()
		if err == io.EOF {
			// log.Fatal(err)
			break
		}
		if err != nil {
			// log.Fatal(err)
			break
		}
		b, err := ioutil.ReadAll(j)
		if err != nil {
			// log.Fatal(err)
			break
		}
		pb := &ds.EntityProto{}
		if err := proto.Unmarshal(b, pb); err != nil {
			log.Fatal(err)
			break
		}
		if onPreload != nil {
			onPreload(dst)
		}
		datastore.LoadEntity(dst, pb)
		onResult(dst)
	}
}