package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/urfave/negroni"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/tobyjsullivan/key-crawler/keys"
	"os"
)

var (
	queueUrl = os.Getenv("KEY_QUEUE_URL")
	received int
	delivered int
)

func main() {
	if queueUrl == "" {
		panic("Must set KEY_QUEUE_URL")
	}

	keyPairs := make(chan *keys.KeyPair, 2000)
	batches := make(chan []*keys.KeyPair, 10)

	client := sqs.New(session.Must(session.NewSession(
		aws.NewConfig().WithCredentials(credentials.NewEnvCredentials()).WithRegion("us-east-1"))))

	r := mux.NewRouter()
	r.HandleFunc("/pairs", pairsPostHandler(keyPairs)).Methods(http.MethodPost)

	n := negroni.New()
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go queueSubmitter(keyPairs, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)
	go sendBatches(client, batches)

	n.Run(":"+port)
}

func pairsPostHandler(keys chan *keys.KeyPair) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		println("Request received.")

		var reqBody pairsReqFmt
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, pair := range reqBody.Pairs {
			keys <- pair
		}
	}
}

func queueSubmitter(keyPairs chan *keys.KeyPair, batches chan []*keys.KeyPair) {
	batchSize := 10
	batch := make([]*keys.KeyPair, 0, batchSize)

	for pair := range keyPairs {
		received++
		if received% 100 == 0 {
			println("address:", pair.Address, "privateKey:", pair.PrivateKey, "received:", received)
		}
		batch = append(batch, pair)

		if len(batch) == batchSize {
			tmp := make([]*keys.KeyPair, batchSize)
			copy(tmp, batch)
			batches<- tmp
			batch = batch[:0]
		}
	}
}

func sendBatches(client *sqs.SQS, batches chan []*keys.KeyPair) {
	for batch := range batches {
		entries := make([]*sqs.SendMessageBatchRequestEntry, 0, len(batch))
		for _, entry := range batch {
			msg, err := json.Marshal(entry)
			if err != nil {
				println("error:", err.Error())
				continue
			}

			entries = append(entries, &sqs.SendMessageBatchRequestEntry{
				Id:          aws.String(entry.Address),
				MessageBody: aws.String(string(msg)),
			})
		}

		_, err := client.SendMessageBatch(&sqs.SendMessageBatchInput{
			QueueUrl: aws.String(queueUrl),
			Entries:  entries,
		})
		if err != nil {
			println("error sending batch:", err.Error())
		}
		delivered += len(batch)
		if delivered%100 == 0 {
			println("delivered:", delivered)
		}
	}
}

type pairsReqFmt struct {
	Pairs []*keys.KeyPair `json:"pairs"`
}
