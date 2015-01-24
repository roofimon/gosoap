package main

import (
	"testing"
)

func TestSanitize(t *testing.T) {
	expected := "string"
	part := Part{Name: "firstName", Type: "xsd:string"}

	part.Sanitize()

	if part.Type != expected {
		t.Errorf("Expected %s but got %s", expected, part.Type)
	}
}

func TestSanitizeMessage(t *testing.T) {
	expected := Part{Name: "firstName", Type: "string"}
	message := Message{Part{Name: "firstName", Type: "xsd:string"}}

	message.Sanitize()

	if message.Part != expected {
		t.Errorf("Expected %v but got %v", expected, message.Part)
	}
}
