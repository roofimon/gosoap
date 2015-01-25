package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

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
}

func (d *Definition) saveToFile() {
	ioutil.WriteFile(d.Name+".go", []byte(d.String()), 0644)
}

type PortType struct {
	Name      string    `xml:"name,attr"`
	Operation Operation `xml:"operation"`
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

type Schema struct {
	ElementFormDefault string `xml:"elementFormDefault,attr"`
	TargetNamespace string `xml:"targetNamespace,attr"`
}

func RemoveNamespace(input string) string {
	return input[strings.Index(input, ":")+1:]
}

var funcMap template.FuncMap = template.FuncMap{
	"title":           strings.Title,
	"removeNamespace": RemoveNamespace,
}

var structTemplate = `package ws
{{range $message := .Messages}}
type {{$message.Name}} struct {
	{{title $message.Part.Name}} {{removeNamespace $message.Part.Type}}
}
{{end}}
func {{title .PortType.Operation.Name}}(req *{{removeNamespace .PortType.Operation.Input.Message}}) (*{{removeNamespace .PortType.Operation.Output.Message}}, error) {
}
`

func (d *Definition) String() string {
	var b bytes.Buffer
	tmpl, _ := template.New("structTemplate").Funcs(funcMap).Parse(structTemplate)
	tmpl.Execute(&b, d)
	return b.String()
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

func main() {
	d := ParseFile("resources/stockquote.wsdl")
	d.Name = "StockQuote"
	fmt.Println(d.Name)
	d.saveToFile()
}
