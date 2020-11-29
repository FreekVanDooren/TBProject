package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"tbp.com/user/hello/history"
	"tbp.com/user/hello/responses"
	"testing"
)

func TestServerUnknownPaths(t *testing.T) {
	server := setupServer()
	defer server.Close()

	for _, path := range []string{
		"/unknown",
		"/primes",
		"/primes/non-prime",
	} {
		t.Run("Does not know "+path, func(t *testing.T) {
			response := doGETRequest(t, server.URL+path)
			defer response.Body.Close()
			if response.StatusCode != 404 {
				t.Errorf("Expected status code 404, but got \"%d\"", response.StatusCode)
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != "404 page not found\n" {
				t.Errorf("Expected body, but got %q", body)
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	server := setupServer()
	defer server.Close()

	primeCases := []struct {
		number   int
		expected responses.Primes
	}{
		{number: 2, expected: responses.Primes{IsPrime: true, Message: "It is prime. Hurray!"}},
		{number: 22, expected: responses.Primes{IsPrime: false, Message: "No"}},
	}

	for _, primeCase := range primeCases {
		t.Run("GETs 200 on integer prime test", func(t *testing.T) {
			response := doGETRequest(t, fmt.Sprintf("%s/primes/%d", server.URL, primeCase.number))
			defer response.Body.Close()

			assertIsPrimeResponse(t, response, primeCase.expected)
		})
	}
}

func TestMessageNotChangeOnRepetitionWithPrime(t *testing.T) {
	server := setupServer(history.Memories{
		23: {10, true},
	})
	defer server.Close()

	t.Run(fmt.Sprintf("Can ask for prime 10 times or more"), func(t *testing.T) {
		response := doGETRequest(t, fmt.Sprintf("%s/primes/%d", server.URL, 23))
		defer response.Body.Close()

		assertIsPrimeResponse(t, response, responses.Primes{IsPrime: true, Message: "It is prime. Hurray!"})
	})
}

func TestMessageChangeOnRepetitionWithNonPrime(t *testing.T) {
	memories := history.Memories{4: {1, false}, 6: {2, false}}
	server := setupServer(memories)
	defer server.Close()

	repeatMessages := []struct {
		number  int
		message string
	}{{4, "No"}, {6, "No, and we already told you so!"}}

	for _, repeatMessage := range repeatMessages {
		t.Run(fmt.Sprintf("Message changes ask for non-prime after %d times", memories[repeatMessage.number].Count+1), func(t *testing.T) {
			response := doGETRequest(t, fmt.Sprintf("%s/primes/%d", server.URL, repeatMessage.number))
			defer response.Body.Close()

			assertIsPrimeResponse(t, response, responses.Primes{IsPrime: false, Message: repeatMessage.message})
		})
	}
}

func TestHistoryEndpoint(t *testing.T) {
	memories := history.Memories{4: {1, false}, 6: {2, false}, 97: {100, true}}
	server := setupServer(memories)
	defer server.Close()

	response := doGETRequest(t, fmt.Sprintf("%s/history/", server.URL))
	defer response.Body.Close()

	assertStatus200(t, response)
	assertJsonHeader(t, response)

	var actual responses.History
	unmarshal(t, response, &actual)

	expected := responses.History{
		Requests: []responses.Request{
			{Number: 4, Count: 1},
			{Number: 6, Count: 2},
			{Number: 97, Count: 100},
		},
	}

	assertFailure := func() {
		t.Errorf("Expected body %+v, but got %+v", expected, actual)
	}

	// Don't know why, but sometimes the comparison of the complete
	// History struct fails with reflect#DeepEqual(), so rolled my own
	if len(actual.Requests) != len(expected.Requests) {
		assertFailure()
	}
	isInExpected := func(request responses.Request) bool {
		for _, expectedRequest := range expected.Requests {
			if expectedRequest == request {
				return true
			}
		}
		return false
	}
	for _, actualRequest := range actual.Requests {
		if !isInExpected(actualRequest) {
			assertFailure()
		}
	}
}

func TestAllowsOnlyDefinedMethods(t *testing.T) {
	server := setupServer()
	defer server.Close()

	testCases := []struct {
		endpoint string
	}{
		{"/history"},
		{"/"},
		{"/primes/23"},
	}
	for _, testCase := range testCases {
		response := doPOSTRequest(t, server.URL+testCase.endpoint, strings.NewReader(""))
		if response.StatusCode != 405 {
			t.Errorf("405 expected, but got %d", response.StatusCode)
		}
	}
}

func TestCurrentResponseMessages(t *testing.T) {
	server := setupServer()
	defer server.Close()

	expected := responses.Messages{
		Messages: []responses.Message{
			{0, "No"},
			{3, "No, and we already told you so!"},
		},
	}

	response, actual := GETMessagesFromServer(t, server)
	defer response.Body.Close()
	if len(actual.Messages) != len(expected.Messages) {
		t.Errorf("Expected body %+v, but got %+v", expected, actual)
	}
}

func TestCanChangeResponseMessages(t *testing.T) {
	server := setupServer(history.Memories{
		22: {8999, false},
		24: {9000, false},
	})
	defer server.Close()

	newMessages := responses.Messages{
		Messages: []responses.Message{
			{0, "No"},
			{3, "No, and we already told you so!"},
			{9001, "It's over 9000!"},
		},
	}

	messageBytes, err := json.Marshal(newMessages)
	if err != nil {
		t.Fatal(err)
	}
	response := doPOSTRequest(t, server.URL+"/messages", bytes.NewReader(messageBytes))
	defer response.Body.Close()

	if response.StatusCode != 202 {
		t.Errorf("Expected status code 202, but got \"%d\"", response.StatusCode)
	}

	response, actualMessages := GETMessagesFromServer(t, server)
	defer response.Body.Close()
	if len(actualMessages.Messages) != len(newMessages.Messages) {
		t.Errorf("Expected body %+v, but got %+v", newMessages, actualMessages)
	}

	response = doGETRequest(t, server.URL+"/primes/22")
	defer response.Body.Close()

	var primes22 responses.Primes
	unmarshal(t, response, &primes22)

	if primes22.Message != "No, and we already told you so!" {
		t.Errorf("Expected message %q, but got %q", "No, and we already told you so!", primes22.Message)
	}

	response = doGETRequest(t, server.URL+"/primes/22")
	defer response.Body.Close()

	var primes24 responses.Primes
	unmarshal(t, response, &primes24)

	if primes24.Message != "It's over 9000!" {
		t.Errorf("Expected message %q, but got %q", "It's over 9000!", primes24.Message)
	}
}

func GETMessagesFromServer(t *testing.T, server *httptest.Server) (*http.Response, responses.Messages) {
	response := doGETRequest(t, fmt.Sprintf("%s/messages/", server.URL))

	assertStatus200(t, response)
	assertJsonHeader(t, response)

	var actual responses.Messages
	unmarshal(t, response, &actual)
	return response, actual
}

func unmarshal(t *testing.T, response *http.Response, actual interface{}) {
	err := json.NewDecoder(response.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}
}

func setupServer(memories ...history.Memories) *httptest.Server {
	if memories == nil {
		return httptest.NewServer(setupRouter(history.Setup()))
	}
	return httptest.NewServer(setupRouter(memories[0]))
}

func doGETRequest(t *testing.T, requestPath string) *http.Response {
	return doRequest(t, requestPath, http.MethodGet, nil)
}

func doPOSTRequest(t *testing.T, requestPath string, body io.Reader) *http.Response {
	return doRequest(t, requestPath, http.MethodPost, body)
}

func doRequest(t *testing.T, requestPath string, method string, body io.Reader) *http.Response {
	request, err := http.NewRequest(method, requestPath, body)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func assertIsPrimeResponse(t *testing.T, response *http.Response, expected responses.Primes) {
	assertStatus200(t, response)
	var actual responses.Primes
	unmarshal(t, response, &actual)

	if actual != expected {
		t.Errorf("Expected body, but got %+v", actual)
	}

	assertJsonHeader(t, response)
}

func assertStatus200(t *testing.T, response *http.Response) {
	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200, but got \"%d\"", response.StatusCode)
	}
}

func assertJsonHeader(t *testing.T, response *http.Response) {
	header := response.Header.Get("Content-Type")
	if header != "application/json" {
		t.Errorf("Expected header, but got %q", header)
	}
}
