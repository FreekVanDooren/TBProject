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
		expected responses.Prime
	}{
		{number: 2, expected: responses.Prime{IsPrime: true, Message: "It is prime. Hurray!"}},
		{number: 22, expected: responses.Prime{IsPrime: false, Message: "No"}},
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

		assertIsPrimeResponse(t, response, responses.Prime{IsPrime: true, Message: "It is prime. Hurray!"})
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

			assertIsPrimeResponse(t, response, responses.Prime{IsPrime: false, Message: repeatMessage.message})
		})
	}
}

func setupServer(memories ...memory.Memories) *httptest.Server {
	if memories == nil {
		return httptest.NewServer(setupRouter(make(memory.Memories)))
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

func assertIsPrimeResponse(t *testing.T, response *http.Response, expected responses.Prime) {
	assertStatus200(t, response)
	var actual responses.Prime
	err := json.NewDecoder(response.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("Expected body, but got %+v", actual)
	}

	header := response.Header.Get("Content-Type")
	if header != "application/json" {
		t.Errorf("Expected header, but got %q", header)
	}
}

func assertStatus200(t *testing.T, response *http.Response) {
	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200, but got \"%d\"", response.StatusCode)
	}
}
