package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/urfave/negroni"
)

func main() {
	keys := make(chan *keyPair, 20)

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

		address := r.FormValue("address")
		privateKey := r.FormValue("private-key")

		if address == "" || privateKey == "" {
			http.Error(w, "must provide address and private-key", http.StatusBadRequest)
			return
		}


		keys <- &keyPair{
			address: address,
			privateKey: privateKey,
		}
	}
}

func queueSubmitter(keys chan *keyPair) {
	for pair := range keys {
		println("address:", pair.address, "privateKey:", pair.privateKey)

		// TODO submit to a queue
	}
}

type keyPair struct {
	address string
	privateKey string
}
