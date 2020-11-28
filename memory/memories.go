package memory

import (
	"tbp.com/user/hello/primes"
	"tbp.com/user/hello/responses"
)

type Memories map[int]*Memory

type Memory struct {
	Count   int
	IsPrime bool
}

func CreateMemories() Memories {
	return make(map[int]*Memory)
}

func (memories *Memories) Update(number int) {
	m := *memories
	var count int
	var isPrime bool
	if m[number] == nil {
		count = 1
		isPrime = primes.IsPrime(number)
	} else {
		count = m[number].Count + 1
		isPrime = m[number].IsPrime
	}
	m[number] = &Memory{Count: count, IsPrime: isPrime}
}

func (memories *Memories) ToPrimeResponse(number int) responses.Primes {
	return (*memories)[number].ToPrimeResponse()
}

func (memories *Memories) ToHistoryResponse() responses.History {
	var requests []responses.Request
	for number, memory := range *memories {
		requests = append(requests, responses.Request{Number: number, Count: memory.Count})
	}
	return responses.History{Requests: requests}
}

func (m *Memory) update() {
	m.Count = m.Count + 1
}

func (m *Memory) toMessage() string {
	if m.IsPrime {
		return "It is prime. Hurray!"
	}
	if m.Count > 2 {
		return "No, and we already told you so!"
	}
	return "No"
}

func (m *Memory) ToPrimeResponse() responses.Primes {
	return responses.Primes{IsPrime: m.IsPrime, Message: m.toMessage()}
}
