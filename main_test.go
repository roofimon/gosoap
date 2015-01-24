package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var message Message = Message{Name: "SayHelloRequest", Part: Part{Name: "firstName", Type: "xsd:string"}}

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
	FirstName string
}

type SayHelloResponse struct {
	Greeting string
}

func SayHello(req *SayHelloRequest) (*SayHelloResponse, error) {
}
`)
	var definition Definition = expectedDefinition

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

func TestRemoveNamespace(t *testing.T) {
	expected := "string"
	input := "tns:string"
	result := RemoveNamespace(input)
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
