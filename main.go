package main

import "strings"

type Part struct {
	Name string
	Type string
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
