// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
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
			if resp, err := db.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(*table),
				Key: map[string]*dynamodb.AttributeValue{
					// "bench_area": &dynamodb.AttributeValue{S: aws.String("KingsLanding")},
					"id": &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%s%d", *idPrefix, rand.Intn(*idTop)))},
				},
			}); err != nil {
				fmt.Println(err)
			} else {
				if *debug {
					log.Println(resp)
				}
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
