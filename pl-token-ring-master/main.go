

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

type Server struct {
	Target chan Token
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method ==  http.MethodPost {
		var token Token
		if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.Target <- token
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	var chainSize = flag.Int("n", 12, "Chain size")
	flag.Parse()
	if chainSize == nil {
		log.Fatal("Expected to provide int value, number of nodes in a token ring")
	}
	tRing := NewTokenRing(*chainSize)
	server := Server{
		Target: tRing.Run(),
	}
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
