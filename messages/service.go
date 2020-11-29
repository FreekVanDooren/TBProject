package messages

import (
	"fmt"
	"log"
	"sort"
	"tbp.com/user/hello/repository"
	"tbp.com/user/hello/responses"
)

type Service struct {
	repository repository.FileRepository
	responses.Messages
}

func (m Service) GetMessage(count int) string {
	messages := m.Messages.Messages
	for _, message := range messages {
		if count >= message.LowerLimit {
			return message.Message
		}
	}
	return messages[len(messages)-1].Message
}

func (m *Service) Update(messages responses.Messages) error {
	err := validate(messages.Messages)
	if err != nil {
		return err
	}
	sort.Sort(messages.Messages)
	m.Messages = messages
	go m.persist()
	return nil
}

func validate(messages responses.MessageSlice) error {
	defaultFound := false
	for _, message := range messages {
		if message.LowerLimit == 0 {
			if defaultFound {
				return fmt.Errorf("must contain only 1 element with lower limit 0, found multiple in %+v", message)
			}
			defaultFound = true
		}
		if message.LowerLimit < 0 {
			return fmt.Errorf("must contain only positive lower limits, found negatives in %+v", message)
		}
	}
	if !defaultFound {
		return fmt.Errorf("must contain element with lower limit 0, found none in %+v", messages)
	}
	return nil
}

func (m Service) Get() interface{} {
	return m.Messages
}

func (m Service) persist() {
	err := m.repository.Persist(m.Messages)
	if err != nil {
		log.Println(err)
	}
}

func Setup(folderName string) (*Service, error) {
	repository, err := repository.Initialize(folderName, "messages")
	var messages responses.Messages
	err = repository.ReadAll(&messages)

	if messages.Messages == nil && err == nil {
		messages = responses.Messages{
			Messages: responses.MessageSlice{
				{3, "No, and we already told you so!"},
				{0, "No"},
			},
		}
		err = repository.Persist(messages)
	}
	return &Service{
		repository: repository,
		Messages:   messages,
	}, err
}
