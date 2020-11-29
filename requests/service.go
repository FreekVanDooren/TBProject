package requests

import (
	"tbp.com/user/hello/messages"
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

func (memories Memories) ToPrimeResponse(number int, messages *messages.FeedbackMessages) responses.Primes {
	return memories[number].ToPrimeResponse(messages)
}

func (memories Memories) ToHistoryResponse() responses.History {
	var requests []responses.Request
	for number, memory := range memories {
		requests = append(requests, responses.Request{Number: number, Count: memory.Count})
	}
	return responses.History{Requests: requests}
}

func (m Memory) update() {
	m.Count = m.Count + 1
}

func (m Memory) toMessage(messages *messages.FeedbackMessages) string {
	if m.IsPrime {
		return "It is prime. Hurray!"
	}
	return messages.GetMessage(m.Count)
}

func (m Memory) ToPrimeResponse(messages *messages.FeedbackMessages) responses.Primes {
	return responses.Primes{IsPrime: m.IsPrime, Message: m.toMessage(messages)}
}
