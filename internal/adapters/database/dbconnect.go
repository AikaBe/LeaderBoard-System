package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := "user=postgres password=yourpassword dbname=yourdb sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}

	log.Println("Подключено к базе данных PostgreSQL")
}
