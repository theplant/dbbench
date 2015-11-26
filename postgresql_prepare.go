// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbURL   = flag.String("url", "postgres://:@localhost/pqgotest?sslmode=disable", "db url")
	total   = flag.Int("t", 100000, "total count to be inserted")
	gophers = flag.Int("g", 8, "gophers count")
	reset   = flag.Bool("reset", false, "reset table")
)

func main() {
	flag.Parse()
	log.Println("Total:", *total)
	log.Println("Gophers:", *gophers)

	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if *reset {
		if _, err := db.Exec("DROP TABLE IF EXISTS test_table"); err != nil {
			panic(err)
		}

		if _, err := db.Exec(`
			CREATE TABLE test_table (
				id         BIGSERIAL NOT NULL PRIMARY KEY,
				email      VARCHAR(255) NOT NULL,
				first_name VARCHAR(255) NOT NULL,
				last_name  VARCHAR(255) NOT NULL,
				age        INTEGER NOT NULL
			)
		`); err != nil {
			panic(err)
		}
	}

	throttle := make(chan int, *gophers)
	start := time.Now()
	var totalDuration int64
	for i := 0; i <= *total; i++ {
		throttle <- i
		if i%(*total/100) == 0 {
			fmt.Printf("\r%s Generated %d%%", time.Now().Format("2006-01-02 15:04:05"), i/(*total/100))
		}
		go func(i int) {
			defer func() { <-throttle }()
			defer func(start time.Time) { atomic.AddInt64(&totalDuration, int64(time.Now().Sub(start))) }(time.Now())
			if _, err := db.Exec(
				"INSERT INTO test_table (email, first_name, last_name, age) VALUES ($1,$2,$3,$4)",
				fmt.Sprintf("laurence-%d@test.com", i),
				fmt.Sprintf("Laurence%d", i),
				"Tester",
				25,
			); err != nil {
				panic(err)
			}
		}(i)
	}
	fmt.Println("")
	log.Println("Took:", time.Now().Sub(start))
	log.Println("Total Duration:", time.Duration(totalDuration))
	log.Println("TPI:", time.Duration(totalDuration/int64(*total)))
}
