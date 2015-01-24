package main

import (
	"reflect"
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

func TestSanitizeDefinition(t *testing.T) {
	expected_message := Message{Part{Name: "firstName", Type: "string"}}
	expected := Definition{Messages: []Message{expected_message}}

	message := Message{Part{Name: "firstName", Type: "xsd:string"}}
	definition := Definition{Messages: []Message{message}}

	definition.Sanitize()

	if !reflect.DeepEqual(definition, expected) {
		t.Errorf(("Expected %v but got %v"), expected, definition)
	}
}

func TestConvertMessage(t *testing.T) {
	expected := Message{
		Part{
			Name: "firstName",
			Type: "xsd:string",
		},
	}

	messageByteArray := []byte(`<message><part name="firstName" type="xsd:string"/></message>`)
	message := ParseWSDLByteArray(messageByteArray)

	if message != expected {
		t.Errorf("Expected %v but got %v", expected, message)
	}
}
