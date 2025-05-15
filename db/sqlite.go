package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var SqliteConn *sql.DB

func Init() {
	var err error
	SqliteConn, err = sql.Open("sqlite", "iot_test.db")
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
		panic(err)
	}
	if err = SqliteConn.Ping(); err != nil {
		SqliteConn.Close()
		log.Fatalf("连接数据库失败: %v", err)
		panic(err)
	}
}
