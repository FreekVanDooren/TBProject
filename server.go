package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"tbp.com/user/hello/history"
	"tbp.com/user/hello/messages"
	"tbp.com/user/hello/responses"
)

func main() {
	memories, err := history.Setup("data")
	if err != nil {
		panic(err)
	}
	feedbackMessages, err := messages.Setup("data")
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", setupRouter(memories, feedbackMessages))
}

func setupRouter(memories history.Service, feedbackMessages *messages.Service) *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("No can do"))
	})
	r.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/history", HistoryHandler(memories)).Methods(http.MethodGet)
	r.HandleFunc("/primes/{number:[0-9]+}", PrimeHandler(memories, feedbackMessages)).Methods(http.MethodGet)
	r.HandleFunc("/messages", GETMessagesHandler(feedbackMessages)).Methods(http.MethodGet)
	r.HandleFunc("/messages", POSTMessagesHandler(feedbackMessages)).Methods(http.MethodPost)
	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Home!")
}

func PrimeHandler(memories history.Service, feedbackMessages *messages.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		potentialNumber := vars["number"]
		number, err := strconv.Atoi(potentialNumber)
		if err != nil {
			fmt.Println(potentialNumber, "is not an integer.")
			http.Error(w, fmt.Sprintf("Not an integer: %s", potentialNumber), http.StatusBadRequest)
			return
		}
		memories.Update(number)
		sendAsJSONResponse(w, memories.ToPrimeResponse(number, feedbackMessages))
	}
}

func HistoryHandler(memories history.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sendAsJSONResponse(w, memories.ToHistoryResponse())
	}
}

func GETMessagesHandler(feedbackMessages *messages.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sendAsJSONResponse(w, feedbackMessages.Get())
	}
}

/*
  Endpoint should have some security... We wouldn't want just anyone updating this.
*/
func POSTMessagesHandler(feedbackMessages *messages.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var messages responses.Messages
		err := json.NewDecoder(r.Body).Decode(&messages)
		//defer r.Body.Close()
		if err != nil {
			if body, err := ioutil.ReadAll(r.Body); err == nil {
				http.Error(w, fmt.Sprintf("Can't unmarshal request from %s", body), http.StatusBadRequest)
			} else {
				http.Error(w, "Can't read body", http.StatusBadRequest)
			}
			return
		}
		err = feedbackMessages.Update(messages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func sendAsJSONResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
