<p align="center" class="img-md mt15 mb15">
	<img src="http://www.sromku.com/static/img/blog/migration-preview.png"/>
</p>

#### Check the blog post: http://www.sromku.com/blog/datastore-sql-migration
---

**Run it yourself**

`go get -u github.com/sromku/datastore-to-sql/backup`

**Load the backup file into model**

This script will load the backup file into 'Profile' model and print it

``` go
import (
	"fmt"
	"github.com/sromku/datastore-to-sql/backup"
)

func main() {
	backupPath := ".../output-0"
	backup.Load(backupPath, &Profile{}, nil,
		func(res interface{}) {
			profile, _ := res.(*Profile) // the loaded model
			fmt.Println(profile)
		})
}

type Profile struct {
	Name   string `datastore:"name, noindex"`
	Email  string `datastore:"email"`
	Gender int `datastore:"gender, noindex"`
}
```