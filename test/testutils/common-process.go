package testutils

import (
	"log"
	"os"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

func PreProcess() {
	setEnv()
	db, err := models.Db()
	if err != nil {
		log.Fatal(err)
	}
	err = db.TruncateTables()
	if err != nil {
		log.Fatal(err)
	}
}

func PostProcess() {
}

func setEnv() error {
	return os.Setenv("MYSQL_CONNECTION_STRING", "root:root@tcp(127.0.0.1:3306)/test?parseTime=true")
}
