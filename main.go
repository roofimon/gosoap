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

func (p *Part) Sanitize() {
	p.Type = strings.Replace(p.Type, "xsd:", "", -1)
}

type Message struct {
	Name string `xml:"name,attr"`
	Part Part   `xml:"part"`
}

func (m *Message) Sanitize() {
	m.Part.Sanitize()
}

type Definition struct {
	Name     string    `xml:"name,attr"`
	Messages []Message `xml:"message"`
	PortType PortType  `xml:"portType"`
}

func (d *Definition) Sanitize() {
	for id := range d.Messages {
		d.Messages[id].Sanitize()
	}
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

var structTemplate = `package ws
{{range $message := .Messages}}
type {{$message.Name}} struct {
	{{$message.Part.Name}} {{$message.Part.Type}}
}
{{end}}`

func (d *Definition) String() string {
	var b bytes.Buffer
	tmpl, _ := template.New("structTemplate").Parse(structTemplate)
	d.Sanitize()
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
