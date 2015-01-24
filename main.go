package main

import (
	"encoding/xml"
	"strings"
)

type Part struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

func (p *Part) Sanitize() {
	p.Type = strings.Replace(p.Type, "xsd:", "", -1)
}

type Message struct {
	Part Part
}

func (m *Message) Sanitize() {
	m.Part.Sanitize()
}

type Definition struct {
	Messages []Message
}

func (d *Definition) Sanitize() {
	for id := range d.Messages {
		d.Messages[id].Sanitize()
	}
}

func ParseWSDLByteArray(partByteArray []byte) Part {
	var part Part
	xml.Unmarshal(partByteArray, &part)
	return part
}
