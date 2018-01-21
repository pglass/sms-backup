package parse

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"time"
)

var (
	CENTRAL_LOCATION *time.Location
)

func init() {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal(err)
	}
	CENTRAL_LOCATION = loc
}

type Document struct {
	XMLName xml.Name `xml:"smses"`
	SMSes   []SMS    `xml:"sms"`
	MMSes   []MMS    `xml:"mms"`
}

func ParseXML(filename string) (Document, error) {
	doc := Document{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return doc, err
	}

	err = xml.Unmarshal(data, &doc)
	if err != nil {
		return doc, err
	}

	return doc, nil
}
