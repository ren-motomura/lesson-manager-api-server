package models

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql" // mysqlに接続するために必要
)

var isDbInitialized = false
var recconectAt time.Time
var dbMap *gorp.DbMap
var mtx sync.Mutex

const recconectPeriod = 300

func Db() (*gorp.DbMap, error) {
	mtx.Lock()
	defer mtx.Unlock()

	if isDbInitialized {
		if recconectAt.After(time.Now()) { // 再接続の必要がない場合
			return dbMap, nil
		}
		if err := dbMap.Db.Close(); err != nil { // 再接続の必要がある場合、接続を一度閉じる
			return dbMap, err
		}
	}

	connectionStr := mustGetenv("MYSQL_CONNECTION_STRING")
	db, err := sql.Open(
		"mysql",
		connectionStr,
	)
	if err != nil {
		return dbMap, err
	}
	recconectAt = time.Now().Add(recconectPeriod)

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

	// register tables here ↓

	registerSession(dbMap)
	registerCompany(dbMap)
	registerStudio(dbMap)
	registerStaff(dbMap)

	// register tables here ↑

	isDbInitialized = true
	return dbMap, nil
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
