package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var initialized = false
var ginLambda *ginadapter.GinLambda
var kc *kinesis.Kinesis
var streamName *string

var (
	streamConfig = flag.String("stream", "super_cool_stream", "your stream name")
	regionConfig = flag.String("region", "us-east-1", "your AWS region")
)

type (
	Data struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}
)

func init() {
	s, _ := session.NewSession(&aws.Config{
		Region:      regionConfig,
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ID"), os.Getenv("AWS_TOKEN"), ""),
	})

	fmt.Println("Creds", os.Getenv("AWS_ID"), os.Getenv("AWS_TOKEN"))

	kc = kinesis.New(s)
	streamName = aws.String(*streamConfig)
}

func main() {
	// get records use shard iterator for making request
	id := "shardId-000000000000"
	// retrieve iterator
	iteratorOutput, _ := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
		ShardId:           &id,
		ShardIteratorType: aws.String("LATEST"),
		// ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
		// ShardIteratorType: aws.String("LATEST"),
		StreamName: streamName,
	})

	records, _ := kc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: iteratorOutput.ShardIterator,
	})
	outputRecords(records.Records)

	for true {
		// and, you can iteratively make GetRecords request using records.NextShardIterator
		recordsSecond, err := kc.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: records.NextShardIterator,
		})
		if err != nil {
			panic(err)
		}
		outputRecords(recordsSecond.Records)
	}
}

func outputRecords(records []*kinesis.Record) {
	for _, record := range records {
		var output Data

		dec := gob.NewDecoder(bytes.NewBuffer(record.Data))
		dec.Decode(&output)

		fmt.Printf("%v\n", output)
	}

}
