package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	db_filename  = "db.sqlite3"
	open_options = "?_fk=true&mode=rwc"
	init_db_stmt = `CREATE TABLE links (
	id INTEGER PRIMARY KEY ASC,
	link TEXT,
	handler TEXT
)`
)

func exec(query string) error {
	db, err := sql.Open("sqlite3", db_filename+open_options)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(query)
	return err
}

func query(query string) (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", db_filename+open_options)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return db.Query(query)
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
