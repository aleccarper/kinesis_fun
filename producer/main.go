package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var initialized = false
var ginLambda *ginadapter.GinLambda
var kc *kinesis.Kinesis
var streamName *string

var (
	streamConfig = flag.String("stream", "super_cool_stream", "your stream name")
	regionConfig = flag.String("region", "us-east-1", "your AWS region")
)

func init() {
	s, _ := session.NewSession(&aws.Config{
		Region:      regionConfig,
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ID"), os.Getenv("AWS_TOKEN"), ""),
	})

	kc = kinesis.New(s)
	streamName = aws.String(*streamConfig)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		r.POST("/produce", doThings)

		ginLambda = ginadapter.New(r)
		initialized = true
	}

	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.Proxy(req)
}

type (
	Data struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}
)

type (
	Payload struct {
		Data []Data `json:"data"`
	}
)

func doThings(c *gin.Context) {
	var json Payload
	if err := c.ShouldBindJSON(&json); err == nil {
		sendData(json.Data)
		c.JSON(http.StatusCreated, nil)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func sendData(data []Data) {
	// put 10 records using PutRecords API
	entries := make([]*kinesis.PutRecordsRequestEntry, len(data))

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	for i := 0; i < len(entries); i++ {
		buffer.Reset()
		enc = gob.NewEncoder(&buffer)
		if err := enc.Encode(data[i]); err != nil {
			log.Fatal(err)
		}
		entries[i] = &kinesis.PutRecordsRequestEntry{
			Data:         buffer.Bytes(),
			PartitionKey: aws.String("1"),
		}
	}

	fmt.Printf("%v\n", entries)
	putsOutput, err := kc.PutRecords(&kinesis.PutRecordsInput{
		Records:    entries,
		StreamName: streamName,
	})
	if err != nil {
		panic(err)
	}

	// putsOutput has Records, and its shard id and sequece enumber.
	fmt.Printf("%v\n", putsOutput)
}

func main() {
	lambda.Start(Handler)
}
