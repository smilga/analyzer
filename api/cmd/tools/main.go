package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
	"github.com/smilga/analyzer/api/datastore/mysql"
)

var migrationsDir = "datastore/mysql/migrations"

func main() {
	operation := flag.String("do", "", "Commands: '-do migrate'")
	flag.Parse()

	db := mysql.NewConnection()
	defer db.Close()

	switch *operation {
	case "migrate":
		migrate(db)
	case "seed":
		seed(db)
	default:
		fmt.Println("No command provided")
	}
}

func migrate(db *sqlx.DB) {
	sqls := readMigrations(migrationsDir)

	for _, sql := range sqls {
		fmt.Printf("Execute: \n\n %s \n\n", sql)
		result, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
		rows, _ := result.RowsAffected()
		fmt.Println(rows)
	}
}

func seed(db *sqlx.DB) {
	userRepo := mysql.NewUserStore(db)
	err := userRepo.Save(admin)
	if err != nil {
		log.Fatal(err)
	}

	patternRepo := mysql.NewPatternStore(db)
	err = patternRepo.Save(isAlive)
	if err != nil {
		log.Fatal(err)
	}
}

func readMigrations(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	migrations := make([]string, len(files))

	for i, file := range files {
		bs, err := ioutil.ReadFile(migrationsDir + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		migrations[i] = string(bs)
	}

	return migrations
}

var admin = &api.User{
	Name:     "admin",
	Email:    "admin@inspected.tech",
	Password: api.Cryptstring("pass"),
	CreatedAt: func() *time.Time {
		now := time.Now()
		return &now
	}(),
}

var isAlive = &api.Pattern{
	Type:        api.System,
	Value:       "isAlive",
	Description: "checks if website is alive",
	CreatedAt: func() *time.Time {
		now := time.Now()
		return &now
	}(),
}
