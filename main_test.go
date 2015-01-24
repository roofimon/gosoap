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

func TestConvertDefinition(t *testing.T) {

	expected := Definition{
		Messages: []Message{
			Message{
				Part: Part{
					Name: "firstName",
					Type: "xsd:string",
				},
			},
		},
	}

	definitionByteArray := []byte(`<definitions name="HelloService"
   targetNamespace="http://www.examples.com/wsdl/HelloService.wsdl"
   xmlns="http://schemas.xmlsoap.org/wsdl/"
   xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/"
   xmlns:tns="http://www.examples.com/wsdl/HelloService.wsdl"
   xmlns:xsd="http://www.w3.org/2001/XMLSchema">
   <message><part name="firstName" type="xsd:string"/></message>
   </definitions>`)

	definition := ParseWSDLByteArray(definitionByteArray)

	if reflect.DeepEqual(expected, definition) {
		t.Errorf("Expected %v but got %v", expected, definition)
	}
}
