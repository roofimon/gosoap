package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
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

var definitionByteArray []byte = []byte(`<definitions name="HelloService"
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
   <portType name="Hello_PortType">
      <operation name="sayHello">
         <input message="tns:SayHelloRequest"/>
         <output message="tns:SayHelloResponse"/>
      </operation>
   </portType>
</definitions>`)

var expectedDefinition Definition = Definition{
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
	PortType: PortType{
		Name: "Hello_PortType",
		Operation: Operation{
			Name:   "sayHello",
			Input:  Input{Message: "tns:SayHelloRequest"},
			Output: Output{Message: "tns:SayHelloResponse"},
		},
	},
}

func TestConvertDefinitionWithMultipleMessages(t *testing.T) {
	definition := ParseWSDLByteArray(definitionByteArray)

	if !reflect.DeepEqual(expectedDefinition, definition) {
		t.Errorf("Expected %v but got %v", expectedDefinition, definition)
	}
}

func TestParseFile(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "")
	ioutil.WriteFile(file.Name(), definitionByteArray, 0777)

	definition := ParseFile(file.Name())

	if !reflect.DeepEqual(expectedDefinition, definition) {
		t.Errorf("Expected %v but got %v", expectedDefinition, definition)
	}
}

func TestWriteFile(t *testing.T) {
	filename := "HelloService.go"
	expected := []byte(`package ws

type SayHelloRequest struct {
	firstName string
}

type SayHelloResponse struct {
	greeting string
}

func sayHello(req *SayHelloRequest) (*SayHelloResponse, error) {
}
`)
	var definition Definition = expectedDefinition
	definition.Sanitize()

	definition.saveToFile()
	data, _ := ioutil.ReadFile(filename)

	if string(data) != string(expected) {
		t.Errorf("Expected %s but got %s", expected, data)
	}

	os.Remove(filename)
}

func TestExtractPortType(t *testing.T) {
	expected := PortType{
		Name: "Hello_PortType",
		Operation: Operation{
			Name:   "sayHello",
			Input:  Input{Message: "tns:SayHelloRequest"},
			Output: Output{Message: "tns:SayHelloResponse"},
		},
	}
	wsdl := []byte(`   <portType name="Hello_PortType">
      <operation name="sayHello">
         <input message="tns:SayHelloRequest"/>
         <output message="tns:SayHelloResponse"/>
      </operation>
   </portType>`)

	var portType PortType
	xml.Unmarshal(wsdl, &portType)

	if expected != portType {
		t.Errorf("Expected %v but got %v", expected, portType)
	}

}

func TestSanitizeInput(t *testing.T) {
	expected := Input{Message: "SayHelloRequest"}
	input := Input{Message: "tns:SayHelloRequest"}

	input.Sanitize()

	if expected != input {
		t.Errorf("Expected %v but got %v", expected, input)
	}
}

func TestSanitizeOutput(t *testing.T) {
	expected := Output{Message: "SayHelloResponse"}
	output := Output{Message: "tns:SayHelloResponse"}

	output.Sanitize()

	if expected != output {
		t.Errorf("Expected %v but got %v", expected, output)
	}
}

func TestSanitizeOperation(t *testing.T) {
	expected := Operation{
		Name:   "sayHello",
		Input:  Input{Message: "SayHelloRequest"},
		Output: Output{Message: "SayHelloResponse"},
	}
	operation := Operation{
		Name:   "sayHello",
		Input:  Input{Message: "tns:SayHelloRequest"},
		Output: Output{Message: "tns:SayHelloResponse"},
	}

	operation.Sanitize()

	if expected != operation {
		t.Errorf("Expected %v but got %v", expected, operation)
	}
}

func TestSanitizePortType(t *testing.T) {
	expected := PortType{
		Operation: Operation{
			Name:   "sayHello",
			Input:  Input{Message: "SayHelloRequest"},
			Output: Output{Message: "SayHelloResponse"},
		},
	}
	portType := PortType{
		Operation: Operation{
			Name:   "sayHello",
			Input:  Input{Message: "tns:SayHelloRequest"},
			Output: Output{Message: "tns:SayHelloResponse"},
		},
	}

	portType.Sanitize()

	if expected != portType {
		t.Errorf("Expected %v but got %v", expected, portType)
	}
}
