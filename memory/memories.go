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
	if (*memories)[number] == nil {
		(*memories)[number] = &Memory{0, primes.IsPrime(number)}
	}
	(*memories)[number].update()
}

func (memories *Memories) ToPrimeResponse(number int) responses.Prime {
	return (*memories)[number].ToPrimeResponse()
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

func (m *Memory) ToPrimeResponse() responses.Prime {
	return responses.Prime{IsPrime: m.IsPrime, Message: m.toMessage()}
}
