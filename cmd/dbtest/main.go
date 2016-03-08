package main

import (
    "log"
    "os"
    "database/sql"
    _ "github.com/lib/pq"
)

var (
    db *sql.DB = nil
)

func main() {

    dburl := os.Getenv("DATABASE_URL")

    log.Printf("DATABASE_URL: %s", dburl)

	db, errd := sql.Open("postgres", dburl)
    if errd != nil {
        log.Fatalf("Error opening database: %q", errd)
    }

    _, err := db.Exec("INSERT INTO twitter_users VALUES ('as', 'aa', 'asdf', 'af')")
	if err != nil {
        log.Fatalf("Error writing database: %q, dburl: %q", err, dburl)
		return
    }

}
