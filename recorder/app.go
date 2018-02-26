package main

import (
	"github.com/tobyjsullivan/key-crawler/keys"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"encoding/json"
	"database/sql"

	_ "github.com/lib/pq"
	"context"
	"os"
)

const (
	queueUrl = "https://sqs.us-east-1.amazonaws.com/110303772622/bitcoin-keys"
)

var (
	received int
	recorded int
	completed int
)

type sqsMessage struct {
	receiptHandle string
	keyPair *keys.KeyPair
}

func main() {
	dbPassword := os.Getenv("PGPASSWORD")
	println("Using DB password:", dbPassword)

	println("recorder starting...")
	client := sqs.New(session.Must(session.NewSession(
		aws.NewConfig().WithCredentials(credentials.NewEnvCredentials()).WithRegion("us-east-1"))))
	db, err := sql.Open("postgres", "")
	if err != nil {
		panic("db connect failed:" + err.Error())
	}

	println("connections established")

	ctx := context.Background()

	inbound := make(chan *sqsMessage, 20)
	completions := make(chan *sqsMessage, 2000)

	go readMessages(client, inbound)

	for i := 0; i < 50; i++ {
		go writeRecords(db, inbound, completions)
	}

	for i := 0; i < 10; i++ {
		go deleteMessages(client, completions)
	}

	<-ctx.Done()
}

func readMessages(client *sqs.SQS, inbound chan *sqsMessage) {
	for {
		out, err := client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueUrl),
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds: aws.Int64(20),
		})
		if err != nil {
			println("error receiving from SQS:", err.Error())
			continue
		}

		for _, msg := range out.Messages {
			handle := *msg.ReceiptHandle
			body := *msg.Body

			var pair keys.KeyPair
			err := json.Unmarshal([]byte(body), &pair)
			if err != nil {
				println("error parsing message body:", err.Error())
				continue
			}

			inbound<- &sqsMessage{
				receiptHandle: handle,
				keyPair: &pair,
			}
			received++
			if received % 100 == 0 {
				println("total received:", received)
			}
		}
	}
}

func writeRecords(db *sql.DB, inbound chan *sqsMessage, completions chan *sqsMessage) {
	for msg := range inbound {
		sqlStatement := "INSERT INTO btckeys (address, private_key) VALUES ($1, $2) ON CONFLICT DO NOTHING"

		_, err := db.Exec(sqlStatement,	msg.keyPair.Address, msg.keyPair.PrivateKey)
		if err != nil {
			println("error writing records to db:", err.Error())
			continue
		}

		completions<- msg
		recorded++
		if recorded % 100 == 0 {
			println("total recorded:", recorded)
		}
	}
}

func deleteMessages(client *sqs.SQS, completions chan *sqsMessage) {
	for msg := range completions {
		_, err := client.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl: aws.String(queueUrl),
			ReceiptHandle: aws.String(msg.receiptHandle),
		})
		if err != nil {
			println("error deleting message:", err.Error())
		}
		completed++
		if completed % 100 == 0 {
			println("total completed:", completed)
		}
	}
}


