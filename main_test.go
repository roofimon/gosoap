package main

import (
	"reflect"
	"testing"
)

var message Message = Message{Name: "SayHelloRequest", Part: Part{Name: "firstName", Type: "xsd:string"}}

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

	message.Sanitize()

	if message.Part != expected {
		t.Errorf("Expected %v but got %v", expected, message.Part)
	}
}

func TestSanitizeDefinition(t *testing.T) {
	expected_message := Message{Name: "SayHelloRequest", Part: Part{Name: "firstName", Type: "string"}}
	expected := Definition{Messages: []Message{expected_message}}

	definition := Definition{Messages: []Message{message}}

	definition.Sanitize()

	if !reflect.DeepEqual(definition, expected) {
		t.Errorf(("Expected %v but got %v"), expected, definition)
	}
}

func TestConvertDefinitionWithMultipleMessages(t *testing.T) {

	expected := Definition{
		Name: "HelloService",
		Messages: []Message{
			Message{
				Name: "SayHelloRequest",
				Part: Part{
					Name: "firstName",
					Type: "xsd:string",
				},
			},
			Message{
				Name: "SayHelloResponse",
				Part: Part{
					Name: "greeting",
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
   <message name="SayHelloRequest">
   	<part name="firstName" type="xsd:string"/>
   	</message>
   <message name="SayHelloResponse">
   	<part name="greeting" type="xsd:string"/>
   	</message>
   </definitions>`)

	definition := ParseWSDLByteArray(definitionByteArray)

	if !reflect.DeepEqual(expected, definition) {
		t.Errorf("Expected %v but got %v", expected, definition)
	}
}
