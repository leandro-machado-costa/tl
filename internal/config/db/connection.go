package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/leandro-machado-costa/tl/internal/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	conf := config.GetDBConfig()

	sc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
		conf.Host, conf.Port, conf.User, conf.Pass, conf.Database)

	var err error
	db, err = sql.Open("postgres", sc)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return err
	}

	log.Println("Database connection opened")
	return nil
}
func GetDB() *sql.DB {
	return db
}
