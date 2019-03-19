package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
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
	default:
		fmt.Println("No command provided")
	}
}

func migrate(db *sqlx.DB) {
	sqls := readMigrations(migrationsDir)

	fmt.Println(sqls)
	for _, sql := range sqls {
		result, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
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
