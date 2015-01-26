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
   <service name="Hello_Service">
      <documentation>WSDL File for HelloService</documentation>
      <port binding="tns:Hello_Binding" name="Hello_Port">
         <soap:address
            location="http://www.examples.com/SayHello/">
      </port>
   </service>
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
	Service: Service {
		Name: "Hello_Service",
		Documentation: "WSDL File for HelloService",
		Port: Port{
			Name: "Hello_Port",
			Binding: "tns:Hello_Binding",
			Address: Address {
				Location: "http://www.examples.com/SayHello/",
			},
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

func TestExtratService (t *testing.T) {
	expected := Service {
		Name: "Hello_Service",
		Documentation: "WSDL File for HelloService",
		Port: Port{
			Name: "Hello_Port",
			Binding: "tns:Hello_Binding",
			Address: Address {
				Location: "http://www.examples.com/SayHello/",
			},
		},
	}

	wsdl := []byte(`   <service name="Hello_Service">
      <documentation>WSDL File for HelloService</documentation>
      <port binding="tns:Hello_Binding" name="Hello_Port">
         <soap:address
            location="http://www.examples.com/SayHello/">
      </port>
   </service>`)

	var service Service
	xml.Unmarshal(wsdl, &service)

	if expected != service {
		t.Errorf("Expected %v but got %v", expected, service)
	}

}

func TestExtractTypes (t *testing.T) {
	expected := Types {
		Schema: Schema {
			ElementFormDefault: "qualified",
			TargetNamespace: "http://www.webserviceX.NET/",
			Elements: []Element {
				Element {
					Name: "GetQuote",
					ComplexType: ComplexType {
						Sequence: Sequence {
							SequenceElement: SequenceElement {
								Name: "symbol",
								Type: "s:string",
							},
						},
					},
				},
				Element {
					Name: "GetQuoteResponse",
					ComplexType: ComplexType {
						Sequence: Sequence {
							SequenceElement: SequenceElement {
								Name: "GetQuoteResult",
								Type: "s:string",
							},
						},
					},
				},
				Element {
					Name: "string",
				},
			},
		},
	}

	wsdl := []byte(`  <wsdl:types>
    <s:schema elementFormDefault="qualified" targetNamespace="http://www.webserviceX.NET/">
      <s:element name="GetQuote">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="symbol" type="s:string" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="GetQuoteResponse">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="GetQuoteResult" type="s:string" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="string" nillable="true" type="s:string" />
    </s:schema>
  </wsdl:types>`)

	var types Types
	xml.Unmarshal(wsdl, &types)

	if !reflect.DeepEqual(expected, types) {
		t.Errorf("Expected %v but got %v", expected, types)
	}

}

func TestTypesWillInDefinition (t *testing.T) {

	wsdl := []byte(`<definitions name="HelloService"
   targetNamespace="http://www.examples.com/wsdl/HelloService.wsdl"
   xmlns="http://schemas.xmlsoap.org/wsdl/"
   xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/"
   xmlns:tns="http://www.examples.com/wsdl/HelloService.wsdl"
   xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <wsdl:types>
 	  <s:schema elementFormDefault="qualified" targetNamespace="http://www.webserviceX.NET/"/>
   </wsdl:types>   
</definitions>`)

	expected := Definition{
		Name: "HelloService",
		Types: Types {
			Schema: Schema {
				ElementFormDefault: "qualified",
				TargetNamespace: "http://www.webserviceX.NET/",
			},
		},
	}

	var definition Definition
	xml.Unmarshal(wsdl, &definition)

	if !reflect.DeepEqual(expected, definition) {
		t.Errorf("Expected %v but got %v", expected, definition)
	}	

}

func TestElementString (t *testing.T) {
	element := Element {
		Name: "GetQuoteResponse",
	}

	expected := `type GetQuoteResponse struct {

}
`
	if expected != element.String() {
		t.Errorf("Expected %s but got %s", expected, element.String())
	}	

}

func TestElementWithUnTitleString (t *testing.T) {
	element := Element {
		Name: "testTitle",
	}

	expected := `type TestTitle struct {

}
`
	if expected != element.String() {
		t.Errorf("Expected %s but got %s", expected, element.String())
	}	

}

func TestPortTypeString (t *testing.T) {
	portType := PortType{
		Name: "Hello_PortType",
		Operation: Operation{
			Name:   "sayHello",
			Input:  Input{Message: "tns:SayHelloRequest"},
			Output: Output{Message: "tns:SayHelloResponse"},
		},
	}


	expected := `func SayHello(req *SayHelloRequest) (*SayHelloResponse, error) {
}
`
	if expected != portType.String() {
		t.Errorf("Expected %s but got %s", expected, portType.String())
	}	

}

func TestStringToTemplateString (t *testing.T) {
	stringTemplate := "func {{.Name}}"
	templateName := "portTypeTemplate"
	structObject := PortType {
		Name: "test",
	}	
	expected := "func test"

	stringResult := StructToTemplateString(templateName, stringTemplate, structObject)

	if expected != stringResult {
		t.Errorf("Expected %s but got %s", expected, stringResult)
	}	
}

func TestStringTypesWithInDefinition (t *testing.T) {
	definition := Definition{
		Name: "HelloService",
		Types: Types {
			Schema: Schema {
				Elements: []Element {
					Element {
						Name: "GetQuote",
					},
					Element {
						Name: "GetQuoteResponse",
					},
				},
			},
		},
	}

	expected := []byte(`package ws

type GetQuote struct {

}

type GetQuoteResponse struct {

}

`)

	data := definition.String()

	if string(data) != string(expected) {
		t.Errorf("Expected %s but got %s", expected, data)
	}	
}