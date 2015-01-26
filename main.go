package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

type TemplateString interface {
    String() string
}

type Part struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Message struct {
	Name string `xml:"name,attr"`
	Part Part   `xml:"part"`
}

type Definition struct {
	Name     string    `xml:"name,attr"`
	Messages []Message `xml:"message"`
	PortType PortType  `xml:"portType"`
	Service  Service   `xml:"service"`
	Types 	 Types     `xml:"types"`
}

func (d *Definition) String() string {

	definitionTemplate := `package ws
{{range $message := .Messages}}
type {{$message.Name}} struct {
	{{title $message.Part.Name}} {{removeNamespace $message.Part.Type}}
}
{{end}}
`	
	returnString := StructToTemplateString("definitionTemplate", definitionTemplate, d)
	returnString += d.Types.String()
	returnString += d.PortType.String()
	return returnString
}

func (d *Definition) saveToFile() error {
	return ioutil.WriteFile(d.Name+".go", []byte(d.String()), 0644)
}

type PortType struct {
	Name      string    `xml:"name,attr"`
	Operation Operation `xml:"operation"`
}

func (p PortType) String() string {

	portTypeTemplate := `{{if .Operation.Name}}func {{title .Operation.Name}}(req *{{removeNamespace .Operation.Input.Message}}) (*{{removeNamespace .Operation.Output.Message}}, error) {
}
{{end}}`
	return StructToTemplateString("portTypeTemplate", portTypeTemplate, p)
}

type Operation struct {
	Name   string `xml:"name,attr"`
	Input  Input  `xml:"input"`
	Output Output `xml:"output"`
}

type Input struct {
	Message string `xml:"message,attr"`
}

type Output struct {
	Message string `xml:"message,attr"`
}

type Service struct {
	Name   string `xml:"name,attr"`
	Documentation string `xml:"documentation"`
	Port Port   `xml:"port"`
}

type Port struct {
	Name   string `xml:"name,attr"`
	Binding   string `xml:"binding,attr"`
	Address Address `xml:"address"`
}

type Address struct {
	Location string `xml:"location,attr"`
}

type Types struct {
	Schema Schema `xml:"schema"`
}

func (t Types) String() string {
	var b bytes.Buffer
	for _, element := range t.Schema.Elements {
		b.WriteString(element.String())
		b.WriteString("\n")
	}

	return b.String()
}

type Schema struct {
	ElementFormDefault string `xml:"elementFormDefault,attr"`
	TargetNamespace string `xml:"targetNamespace,attr"`
	Elements []Element `xml:"element"`
}

type Element struct {
	Name string `xml:"name,attr"`
	ComplexType ComplexType `xml:"complexType"`
}

func (e Element) String() string {

	elementTemplate := `type {{title .Name}} struct {

}
`
	return StructToTemplateString("elementTemplate", elementTemplate, e)
}

type ComplexType struct {
	Sequence Sequence `xml:"sequence"`
}

type Sequence struct {
	SequenceElement SequenceElement `xml:"element"`
}

type SequenceElement struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

var funcMap template.FuncMap = template.FuncMap{
	"title":           strings.Title,
	"removeNamespace": RemoveNamespace,
}

func ParseWSDLByteArray(definitionByteArray []byte) Definition {
	var definition Definition
	xml.Unmarshal(definitionByteArray, &definition)
	return definition
}

func ParseFile(filename string) Definition {
	data, _ := ioutil.ReadFile(filename)
	return ParseWSDLByteArray(data)
}

func RemoveNamespace(input string) string {
	return input[strings.Index(input, ":")+1:]
}

func StructToTemplateString(templateName string, templateString string, structTemplate TemplateString) string {
	var b bytes.Buffer
	tmpl, _ := template.New(templateName).Funcs(funcMap).Parse(templateString)
	tmpl.Execute(&b, structTemplate)
	return b.String()
}

func main() {
	d := ParseFile("resources/stockquote.wsdl")
	d.Name = "StockQuote"
	fmt.Println(d.Name)
	d.saveToFile()
}
