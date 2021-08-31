package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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

type Body struct {
	Url     string `binding:"required"`
	Handler string
}

func exec(query string, args ...interface{}) error {
	db, err := sql.Open("sqlite3", db_filename+open_options)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(query, args...)
	return err
}

func query(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", db_filename+open_options)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return db.Query(query, args...)
}

func insert_link(link, handler string) (string, error) {
	if handler == "" {
		hash := md5.Sum([]byte(link))
		handler = hex.EncodeToString(hash[:])[:6]
	}
	err := exec(`INSERT INTO links (link, handler) VALUES (?, ?)`, link, handler)
	return handler, err
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
		} else {
			c.String(404, "No link found for "+hash)
		}
	})
	r.POST("/link", func(c *gin.Context) {
		var body Body
		c.BindJSON(&body)
		handler, err := insert_link(body.Url, body.Handler)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"handler": handler})
	})
	r.NoRoute(func(c *gin.Context) {
		c.String(404, "Not found")
	})
	r.Run()
}
