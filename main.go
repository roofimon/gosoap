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
	p.Name = strings.ToUpper(string(p.Name[0])) + p.Name[1:]
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
	d.PortType.Sanitize()
}

func (d *Definition) saveToFile() {
	ioutil.WriteFile(d.Name+".go", []byte(d.String()), 0644)
}

type PortType struct {
	Name      string    `xml:"name,attr"`
	Operation Operation `xml:"operation"`
}

func (p *PortType) Sanitize() {
	p.Operation.Sanitize()
}

type Operation struct {
	Name   string `xml:"name,attr"`
	Input  Input  `xml:"input"`
	Output Output `xml:"output"`
}

func (o *Operation) Sanitize() {
	o.Input.Sanitize()
	o.Output.Sanitize()
}

type Input struct {
	Message string `xml:"message,attr"`
}

func (i *Input) Sanitize() {
	i.Message = strings.Replace(i.Message, "tns:", "", -1)
}

type Output struct {
	Message string `xml:"message,attr"`
}

func (i *Output) Sanitize() {
	i.Message = strings.Replace(i.Message, "tns:", "", -1)
}

var structTemplate = `package ws
{{range $message := .Messages}}
type {{$message.Name}} struct {
	{{$message.Part.Name}} {{$message.Part.Type}}
}
{{end}}
func {{.PortType.Operation.Name}}(req *{{.PortType.Operation.Input.Message}}) (*{{.PortType.Operation.Output.Message}}, error) {
}
`

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
