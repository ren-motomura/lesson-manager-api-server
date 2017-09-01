package testutils

import (
	"log"
	"os"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

type gorpLogger struct{}

func (gl *gorpLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func PreProcess() {
	setEnv()
	db, err := models.Db()
	if err != nil {
		log.Fatal(err)
	}
	db.TraceOn("[gorp SQL log]", &gorpLogger{})
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
