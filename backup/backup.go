package backup

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"io/ioutil"
	"log"
	"os"
	pb "github.com/sromku/datastore-to-sql/backup/pb"
)

// Load the backup into real model
// backupFilePath - the backup file path
// dst - the struct model that represents datastore entity and the model you want to load the data of this backup
// onPreload - callback that will be called before loading each entity
// onResult - callback that will be called with already loaded entity in the model
func Load(backupFilePath string, dst interface{}, onPreload func(dst interface{}), onResult func(dst interface{})) {
	f, err := os.Open(backupFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	journals := journal.NewReader(f, nil, false, true)
	for {
		j, err := journals.Next()
		if err != nil {
			// log.Fatal(err)
			break
		}
		b, err := ioutil.ReadAll(j)
		if err != nil {
			// log.Fatal(err)
			break
		}
		pb := &pb.EntityProto{}
		if err := proto.Unmarshal(b, pb); err != nil {
			log.Fatal(err)
			break
		}
		if onPreload != nil {
			onPreload(dst)
		}
		LoadEntity(dst, pb)
		onResult(dst)
	}
}