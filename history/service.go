package history

import (
	"log"
	"tbp.com/user/hello/messages"
	"tbp.com/user/hello/repository"
	"tbp.com/user/hello/responses"
)

type Service struct {
	repository repository.FileRepository
	memories   Memories
}

func Setup(folderName string) (Service, error) {
	repository, err := repository.Initialize(folderName, "history")
	var memories Memories
	err = repository.ReadAll(&memories)
	if memories == nil && err == nil {
		memories = make(map[int]*Memory)
		err = repository.Persist(memories)
	}
	service := Service{memories: memories, repository: repository}
	return service, err
}

func SetupWith(memories Memories, folderName string) (Service, error) {
	repository, err := repository.Initialize(folderName, "history")
	return Service{memories: memories, repository: repository}, err
}

func (s Service) ToHistoryResponse() responses.History {
	return s.memories.ToHistoryResponse()
}

func (s Service) Update(number int) {
	s.memories.Update(number)
	go s.persist()
}

func (s Service) ToPrimeResponse(number int, feedbackMessages *messages.Service) responses.Primes {
	return s.memories.ToPrimeResponse(number, feedbackMessages)
}

func (s Service) persist() {
	err := s.repository.Persist(s.memories)
	if err != nil {
		log.Println(err)
	}
}
