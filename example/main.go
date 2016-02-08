package main

import (
	"fmt"
	"github.com/sromku/datastore-to-sql/backup"
)

func main() {

	backupPaths := []string{
		"../exported-data/data/datastore_backup_datastore_backup_2016_02_07_Profile/157249940434231075281045461947F/output-0",
	}

	for _, backupPath := range backupPaths {

		backup.Load(backupPath, &Profile{},

			// if you want to prepare something before loading into model
			func(res interface{}) {
				profile, _ := res.(*Profile)
				profile.Name = ""
				profile.Email = ""
				profile.Gender = 0
			},

			// process the loaded model
			func(res interface{}) {
				profile, _ := res.(*Profile)
				insert := "INSERT INTO `users` (name,email,gender) VALUES ('%v','%v',%v);"
				insert = fmt.Sprintf(insert, profile.Name, profile.Email, profile.Gender)
				fmt.Println(insert)
			})
	}

}
