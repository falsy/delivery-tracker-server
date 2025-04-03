package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	log.Println("DB 연결 시도")

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB 핑 실패:", err)
	}

	log.Println("DB 연결")
}
