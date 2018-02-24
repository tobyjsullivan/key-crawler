package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/urfave/negroni"
	"encoding/json"
)

func main() {
	keys := make(chan *keyPair, 2000)

	r := mux.NewRouter()
	r.HandleFunc("/pairs", pairsPostHandler(keys)).Methods(http.MethodPost)

	n := negroni.New()
	n.UseHandler(r)

	port := "3000"

	go queueSubmitter(keys)

	n.Run(":"+port)
}

func pairsPostHandler(keys chan *keyPair) func(http.ResponseWriter, *http.Request) {
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

func queueSubmitter(keys chan *keyPair) {
	var count int
	for pair := range keys {
		count++
		println("address:", pair.Address, "privateKey:", pair.PrivateKey, "received:", count)

		// TODO submit to a queue
	}
}

type keyPair struct {
	Address string `json:"address"`
	PrivateKey string `json:"private-key"`
}

type pairsReqFmt struct {
	Pairs []*keyPair `json:"pairs"`
}
