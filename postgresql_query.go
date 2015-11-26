// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbURL   = flag.String("url", "postgres://:@localhost/pqgotest?sslmode=disable", "db url")
	total   = flag.Int("t", 100000, "total count to be inserted")
	idTop   = flag.Int("id", 100000, "id top")
	gophers = flag.Int("g", 16, "concurrent gopher requests")
	verbose = flag.Bool("v", false, "verbose")
	debug   = flag.Bool("d", false, "debug")
)

func main() {
	flag.Parse()
	log.Println("Total:", *total)
	log.Println("Gophers:", *gophers)

	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatal(err)
	}

	jobChans := make(chan struct{}, *gophers)
	start := time.Now()
	var totalDuration int64
	var wg sync.WaitGroup

	rand.Seed(time.Now().Unix())
	for i := 0; i <= *total; i++ {
		jobChans <- struct{}{}
		wg.Add(1)
		if i%((*total)/100) == 0 {
			fmt.Printf("\r%s Queried %d%%", time.Now().Format("2006-01-02 15:04:05"), i/((*total)/100))
		}
		go func() {
			start := time.Now()
			rand.Intn(*idTop)
			var id, age int
			var email, firstName, lastName string
			if err := db.QueryRow("select * from test_table where id = $1", rand.Intn(*idTop)).Scan(&id, &email, &firstName, &lastName, &age); err != nil {
				panic(err)
			}
			if *debug {
				log.Println(id, email, firstName, lastName, age)
			}

			atomic.AddInt64(&totalDuration, int64(time.Now().Sub(start)))
			wg.Done()
			<-jobChans
		}()
	}

	wg.Wait()

	fmt.Println("")
	log.Println("Took:", time.Now().Sub(start))
	log.Println("Total Duration:", time.Duration(totalDuration))
	log.Println("TPQ:", time.Duration(totalDuration/int64(*total)))
}
