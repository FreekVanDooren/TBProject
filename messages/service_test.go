package messages

import (
	"tbp.com/user/hello/responses"
	"testing"
)

func TestUpdatesOnlyWith0Present(t *testing.T) {
	service, err := Setup("test_data")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Updates without error when one default element present", func(t *testing.T) {
		messages := responses.Messages{
			Messages: []responses.Message{
				{0, "No"},
				{3, "No, and we already told you so!"},
			},
		}
		if service.Update(messages) != nil {
			t.Error("Expected no error")
		}
	})
	t.Run("Can update only when default element present", func(t *testing.T) {
		messages := responses.Messages{
			Messages: []responses.Message{
				{3, "No, and we already told you so!"},
			},
		}
		if service.Update(messages) == nil {
			t.Error("Expected an error")
		}
	})
	t.Run("Can update only when single default element present", func(t *testing.T) {
		messages := responses.Messages{
			Messages: []responses.Message{
				{0, "No"},
				{0, "No"},
				{3, "No, and we already told you so!"},
			},
		}
		if service.Update(messages) == nil {
			t.Error("Expected an error")
		}
	})
	t.Run("Can update only when all lower limits are positive", func(t *testing.T) {
		messages := responses.Messages{
			Messages: []responses.Message{
				{0, "No"},
				{-1, "No, please"},
				{3, "No, and we already told you so!"},
			},
		}
		if service.Update(messages) == nil {
			t.Error("Expected an error")
		}
	})
}
