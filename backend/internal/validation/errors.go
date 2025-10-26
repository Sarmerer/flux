package validation

import (
	"fmt"
	"strings"
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func (e Errors) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}
