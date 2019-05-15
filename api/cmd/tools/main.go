package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
	"github.com/smilga/analyzer/api/datastore/cache"
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
	case "uadd":
		createUser(db)
	case "seed":
		seed(db)
	default:
		fmt.Println("No command provided")
	}
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&api.Cookie{},
		&api.Script{},
		&api.Link{},
		&api.Header{},
		&api.HTMLSource{},
		&api.Result{},
	)
	// sqls := readMigrations(migrationsDir)

	// for _, sql := range sqls {
	// 	fmt.Printf("Execute: \n\n %s \n\n", sql)
	// 	result, err := db.Exec(sql)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	rows, _ := result.RowsAffected()
	// 	fmt.Println(rows)
	// }
}

func seed(db *gorm.DB) {
	userRepo := mysql.NewUserStore(db)
	err := userRepo.Save(admin)
	if err != nil {
		log.Fatal(err)
	}

	patternRepo := cache.NewPatternCache(mysql.NewPatternStore(db.DB()))
	err = patternRepo.Save(isAlive)
	if err != nil {
		log.Fatal(err)
	}

	err = patternRepo.Save(hasError)
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

func createUser(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nEnter password: ")
	pass, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	user := &api.User{
		Name:     strings.TrimSpace(email),
		Email:    strings.TrimSpace(email),
		Password: api.Cryptstring(strings.TrimSpace(pass)),
		CreatedAt: func() *time.Time {
			now := time.Now()
			return &now
		}(),
	}

	userRepo := mysql.NewUserStore(db)
	err = userRepo.Save(user)
	if err != nil {
		log.Fatal(err)
	}

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

var hasError = &api.Pattern{
	Type:        api.System,
	Value:       "hasError",
	Description: "Error crawling page",
	CreatedAt: func() *time.Time {
		now := time.Now()
		return &now
	}(),
}
