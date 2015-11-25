// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	idPrefix = flag.String("prefix", "", "id prefix")
	table    = flag.String("table", "test_table", "table name")
	total    = flag.Int("t", 100000, "total count to be inserted")
	idTop    = flag.Int("id", 100000, "id top")
	gophers  = flag.Int("g", 16, "concurrent gopher requests")
	verbose  = flag.Bool("v", false, "verbose")
	debug    = flag.Bool("d", false, "debug")
	endpoint = flag.String("endpoint", "http://localhost:8000", "endpoint")
	region   = flag.String("region", "us-west-2", "region")
)

func main() {
	flag.Parse()
	log.Println("Total:", *total)
	log.Println("Gophers:", *gophers)
	log.Println("ID Top:", *idTop)
	if *debug {
		*verbose = true
	}

	db := dynamodb.New(session.New(
		aws.NewConfig().
			WithEndpoint(*endpoint).
			WithRegion(*region).
			WithCredentials(credentials.NewEnvCredentials()),
	))

	jobChans := make(chan struct{}, *gophers)
	// var wg sync.WaitGroup
	var count int
	var lock sync.Mutex
	var totalDuration time.Duration
	doneChan := make(chan int)
	start := time.Now()

	rand.Seed(time.Now().Unix())
	go func() {
		for {
			jobChans <- struct{}{}
			go func() {
				start := time.Now()
				if resp, err := db.GetItem(&dynamodb.GetItemInput{
					TableName: aws.String(*table),
					Key: map[string]*dynamodb.AttributeValue{
						// "bench_area": &dynamodb.AttributeValue{S: aws.String("KingsLanding")},
						"id": &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%s%d", *idPrefix, rand.Intn(*idTop)))},
					},
				}); err != nil {
					panic(err)
				} else {
					if *debug {
						log.Println(resp)
					}
				}

				lock.Lock()
				count++
				if count%((*total)/100) == 0 {
					if *verbose {
						fmt.Printf("\r%s Queried %d%%", time.Now().Format("2006-01-02 15:04:05"), count/((*total)/100))
					}
				}
				totalDuration += time.Now().Sub(start)
				if count >= *total {
					doneChan <- count
				}
				lock.Unlock()
				<-jobChans
			}()
		}
	}()

	fcount := <-doneChan
	if *verbose {
		fmt.Println("")
	}
	log.Println("Count:", fcount)
	log.Println("Took:", totalDuration)
	log.Println("PQT:", totalDuration/time.Duration(*total))
	log.Println("Took:", time.Now().Sub(start))
}
