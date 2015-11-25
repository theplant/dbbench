// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	idPrefix = flag.String("prefix", "", "id prefix")
	table    = flag.String("table", "test_table", "table name")
	total    = flag.Int("t", 100000, "total count to be inserted")
	gophers  = flag.Int("g", 8, "gophers count")
	reset    = flag.Bool("reset", false, "reset table")
	endpoint = flag.String("endpoint", "http://localhost:8000", "endpoint")
	region   = flag.String("region", "us-west-2", "region")
)

func main() {
	flag.Parse()
	log.Println("Total:", *total)
	log.Println("Table:", *table)

	db := dynamodb.New(session.New(
		aws.NewConfig().
			WithEndpoint(*endpoint).
			WithRegion(*region).
			WithCredentials(credentials.NewEnvCredentials()),
	))

	if *reset {
		if resp, err := db.ListTables(&dynamodb.ListTablesInput{}); err != nil {
			panic(err)
		} else {
			for _, name := range resp.TableNames {
				if *name == *table {
					log.Println("DeleteTable")
					if _, err := db.DeleteTable(&dynamodb.DeleteTableInput{
						TableName: aws.String(*table),
					}); err != nil {
						panic(err)
					}
					break
				}
			}
		}

		log.Println("CreateTable")
		if _, err := db.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(*table),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				// &dynamodb.AttributeDefinition{
				// 	AttributeName: aws.String("bench_area"),
				// 	AttributeType: aws.String("S"),
				// },
				&dynamodb.AttributeDefinition{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("S"),
				},
				// &dynamodb.AttributeDefinition{
				// 	AttributeName: aws.String("email"),
				// 	AttributeType: aws.String("S"),
				// },
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				// &dynamodb.KeySchemaElement{
				// 	AttributeName: aws.String("bench_area"),
				// 	KeyType:       aws.String("HASH"),
				// },
				// &dynamodb.KeySchemaElement{
				// 	AttributeName: aws.String("id"),
				// 	KeyType:       aws.String("RANGE"),
				// },
				&dynamodb.KeySchemaElement{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
			// GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			// 	&dynamodb.GlobalSecondaryIndex{
			// 		IndexName: aws.String("email-index"),
			// 		KeySchema: []*dynamodb.KeySchemaElement{
			// 			&dynamodb.KeySchemaElement{
			// 				AttributeName: aws.String("email"),
			// 				KeyType:       aws.String("HASH"),
			// 			},
			// 		},
			// 		Projection: &dynamodb.Projection{
			// 			// NonKeyAttributes: []*string{aws.String("email")},
			// 			ProjectionType: aws.String("KEYS_ONLY"),
			// 		},
			// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			// 			ReadCapacityUnits:  aws.Int64(1),
			// 			WriteCapacityUnits: aws.Int64(1),
			// 		},
			// 	},
			// },
		}); err != nil {
			panic(err)
		}
	}

	throttle := make(chan int, *gophers)
	start := time.Now()
	var totalDuration time.Duration
	for i := 0; i <= *total; i++ {
		throttle <- i
		if i%(*total/100) == 0 {
			fmt.Printf("\r%s Generated %d%%", time.Now().Format("2006-01-02 15:04:05"), i/(*total/100))
		}
		go func(i int) {
			defer func() { <-throttle }()
			defer func(start time.Time) { totalDuration += time.Now().Sub(start) }(time.Now())
			if _, err := db.PutItem(&dynamodb.PutItemInput{
				TableName: aws.String(*table),
				Item: map[string]*dynamodb.AttributeValue{
					"bench_area": &dynamodb.AttributeValue{S: aws.String("KingsLanding")},
					"id":         &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("%s%d", *idPrefix, i))},
					"email":      &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("test-%d@test.com", i))},
					"first_name": &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("Laurence-%d", i))},
					"last_name":  &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("Tester-%d", i))},
					"age":        &dynamodb.AttributeValue{N: aws.String("25")},
				},
			}); err != nil {
				panic(err)
			}
		}(i)
	}
	fmt.Println("")
	log.Println("Total Duration:", totalDuration)
	log.Println("TPI:", totalDuration/time.Duration(*total))
	log.Println("Took:", time.Now().Sub(start))
}
