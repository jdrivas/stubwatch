package hublib

import (
  "encoding/xml"
)

type StubhubError struct {
  XMLName xml.Name    `xml:"Envelope"`
  Body SHBody
}

type SHBody struct {
  XMLName xml.Name `xml:"Body"`
  Fault SHFault
}

type SHFault struct {
  XMLName xml.Name    `xml:"fault"`
  Code string         `xml:"code"`
  Message string      `xml:"message"`
  Description string  `xml:"description"`
}

func (e StubhubError) Fault() (SHFault) {
  return e.Body.Fault
}

