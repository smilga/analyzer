package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func NewConnection() *gorm.DB {
	tries := 10
	var database *gorm.DB
	var err error

	for tries > 0 {
		fmt.Println("Trying to connect database tries: ", tries)
		if db, e := connect(); e != nil {
			fmt.Println("Error connecting db, ", err)
			time.Sleep(time.Second * 10)
			err = e
			tries--
			continue
		} else {
			tries = 0
			database = db
		}
	}

	if database != nil {
		return database
	}

	log.Fatal("Cant connect to database", err)
	return nil
}

func connect() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	db := os.Getenv("MYSQL_DATABASE")

	// DB, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(database:3306)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", user, pass, db))
	DB, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(database:3306)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", user, pass, db))
	if err != nil {
		return nil, err
	}

	return DB, nil
}
