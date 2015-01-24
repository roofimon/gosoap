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
