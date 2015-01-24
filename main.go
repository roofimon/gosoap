package main

import (
	"bytes"
	"encoding/xml"
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
