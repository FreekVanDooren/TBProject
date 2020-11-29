package messages

import (
	"fmt"
	"sort"
	"tbp.com/user/hello/responses"
)

type FeedbackMessages struct {
	responses.Messages
}

func (m FeedbackMessages) GetMessage(count int) string {
	messages := m.Messages.Messages
	for _, message := range messages {
		if count >= message.LowerLimit {
			return message.Message
		}
	}
	return messages[len(messages)-1].Message
}

func (m *FeedbackMessages) Update(messages responses.Messages) error {
	err := validate(messages.Messages)
	if err != nil {
		return err
	}
	sort.Sort(messages.Messages)
	m.Messages = messages
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

func (m FeedbackMessages) Get() interface{} {
	return m.Messages
}

func CreateMessages() *FeedbackMessages {
	messages := responses.Messages{
		Messages: []responses.Message{
			{3, "No, and we already told you so!"},
			{0, "No"},
		},
	}
	return &FeedbackMessages{Messages: messages}
}
