package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	dsn := "root:wujinlei@tcp(127.0.0.1:3306)/task_schedule_center"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("db connect success...")
}
