package main

import (
	"os"
	"strconv"
	"log"
	"time"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"encoding/json"
)

const (
	loopDelay = 1 * time.Second
)

func main() {
	current := requireEnvInt("START", 64)
	batchSize := requireEnvInt("BATCH_SIZE", 32)
	sqsQueueUrl := requireEnvString("BATCH_QUEUE_URL")

	client := sqs.New(session.Must(session.NewSession(
		aws.NewConfig().WithCredentials(credentials.NewEnvCredentials()).WithRegion("us-east-1"))))

	batchQueue := make(chan int64)

	go submitBatches(client, batchQueue, sqsQueueUrl, batchSize)

	for current > 0 {
		batchQueue <- current

		current += batchSize
	}
}

func requireEnvString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("[requireEnvString]", "Error: not specified:", key)
	}

	return val
}

func requireEnvInt(key string, bitSize int) int64 {
	strVal := requireEnvString(key)

	val, err := strconv.ParseInt(strVal, 10, bitSize)
	if err != nil {
		log.Fatal("[requireEnvInt]", "Error: parse failed:", err)
	}

	return val
}

type fmtBatch struct {
	Start int64 `json:"start"`
	Size  int64 `json:"size"`
}

func submitBatches(client *sqs.SQS, batches chan int64, queueUrl string, batchSize int64) {
	for start := range batches {
		batch := &fmtBatch{
			Start: start,
			Size:  batchSize,
		}

		msg, err := json.Marshal(batch)
		if err != nil {
			log.Println("[submitBatches] Error: failed to marshal batch:", err.Error())
			retry(start, batches)
			continue
		}

		log.Println("[submitBatches]", "submitting:", string(msg))

		_, err = client.SendMessage(&sqs.SendMessageInput{
			QueueUrl:               aws.String(queueUrl),
			MessageBody:            aws.String(string(msg)),
		})
		if err != nil {
			log.Println("[submitBatches] Error: SendMessage failed:", err.Error())
			retry(start, batches)
			continue
		}
	}
}

func retry(start int64, batches chan int64) {
	go func() {
		batches <- start
	}()
}
