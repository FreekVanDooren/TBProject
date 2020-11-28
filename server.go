package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tbp.com/user/hello/memory"
)

func main() {
	http.ListenAndServe(":8080", setupRouter(memory.CreateMemories()))
}

func setupRouter(memories memory.Memories) *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("No can do"))
	})
	r.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/primes/{number:[0-9]+}", PrimeHandler(memories)).Methods(http.MethodGet)
	r.HandleFunc("/history", HistoryHandler(memories)).Methods(http.MethodGet)
	return r
}

func HistoryHandler(memories memory.Memories) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		respondAsJSON(w, memories.ToHistoryResponse())
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Home!")
}

func PrimeHandler(memories memory.Memories) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		potentialNumber := vars["number"]
		number, err := strconv.Atoi(potentialNumber)
		if err != nil {
			fmt.Println(potentialNumber, "is not an integer.")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Not an integer: %s", potentialNumber)
			return
		}
		memories.Update(number)
		respondAsJSON(w, memories.ToPrimeResponse(number))
	}
}

func respondAsJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
