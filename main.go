package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

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
	prune_db_on_launch := flag.Bool("prune", false, "prune database on launch")
	flag.Parse()
	if *prune_db_on_launch {
		log.Println("Pruning database...")
		os.Remove(db_filename)
	}

	if _, err := os.Stat(db_filename); os.IsNotExist(err) {
		log.Println("Creating database...")
		if err := exec(init_db_stmt); err != nil {
			panic(err)
		}
	}
	r := gin.Default()
	r.GET("/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		rows, err := query("SELECT link FROM links WHERE handler = '" + hash + "'")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		if rows.Next() {
			var link string
			if err := rows.Scan(&link); err != nil {
				panic(err)
			}
			c.Redirect(302, link)
		}
		c.String(404, "No link found for "+hash)
	})
	r.NoRoute(func(c *gin.Context) {
		c.String(404, "Not found")
	})
	r.Run()
}
