package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	// Using postgres sql driver
	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

var (
	// DB returns a gorm.DB interface, it is used to access to database
	DB *gorm.DB
)

type Answer struct {
	Question string   `json:"question" binding:"required"`
	Answers  []string `json:"answers" binding:"required"`
}

type Answers []Answer

func (as Answers) Value() (driver.Value, error) {
	return json.Marshal(as)
}

func (as *Answers) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), as)
}

// Survey stores all interview answers
type Survey struct {
	gorm.Model
	Answers Answers `sql:"type:json;not null;default:'[]'" json:"answers" binding:"required"`
}

func init() {
	initDB()
	migrate()
}

func initDB() {
	var err error
	var db gorm.DB

	dbParams := os.Getenv("DB_PARAMS")
	if dbParams == "" {
		panic(errors.New("DB_PARAMS environment variable not set"))
	}

	db, err = gorm.Open("postgres", fmt.Sprintf(dbParams))
	if err == nil {
		DB = &db
	} else {
		panic(err)
	}
}

func migrate() {
	DB.DropTableIfExists(&Survey{})
	DB.AutoMigrate(&Survey{})
}
