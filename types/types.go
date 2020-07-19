package types

import "encoding/xml"

type Gather struct {
	XMLName             xml.Name `xml:"Gather"`
	Action              string   `xml:"action,attr,omitempty"`
	Method              string   `xml:"method,attr,omitempty"`
	NumDigits           string   `xml:"numDigits,attr,omitempty"`
	Timeout             string   `xml:"timeout,attr,omitempty"`
	ActionOnEmptyResult string   `xml:"actionOnEmptyResult,attr,omitempty"`
	Pause               string   `xml:",omitempty"`
	Say                 string   `xml:",omitempty"`
	Play                *Play    `xml:",omitempty"`
}

type Play struct {
	XMLName xml.Name `xml:"Play"`
	Digits  string   `xml:"digits,attr,omitempty"`
	Message string   `xml:",chardata"`
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`
	Say     string   `xml:",omitempty"`
	Play    *Play    `xml:",omitempty"`
	Gather  *Gather  `xml:",omitempty"`
}
