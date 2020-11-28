package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"tbp.com/user/hello/memory"
	"tbp.com/user/hello/responses"
	"testing"
)

func TestServerUnknownPaths(t *testing.T) {
	server := setupServer()
	defer server.Close()

	for _, path := range []string{
		"/unknown",
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
	server := setupServer(memory.Memories{
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
	memories := memory.Memories{4: {1, false}, 6: {2, false}}
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

func TestHistoryHandler(t *testing.T) {
	memories := memory.Memories{4: {1, false}, 6: {2, false}, 97: {100, true}}
	server := setupServer(memories)
	defer server.Close()

	response := doGETRequest(t, fmt.Sprintf("%s/history/", server.URL))
	defer response.Body.Close()

	assertStatus200(t, response)
	assertJsonHeader(t, response)

	var actual responses.History
	err := json.NewDecoder(response.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

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

	// Don't know why, but sometimes the comparison of the complete History struct fails, so rolled my own
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

func setupServer(memories ...memory.Memories) *httptest.Server {
	if memories == nil {
		return httptest.NewServer(setupRouter(memory.CreateMemories()))
	}
	return httptest.NewServer(setupRouter(memories[0]))
}

func doGETRequest(t *testing.T, requestPath string) *http.Response {
	request, err := http.NewRequest("GET", requestPath, nil)
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
	err := json.NewDecoder(response.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

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
