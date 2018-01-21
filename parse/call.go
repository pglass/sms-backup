package parse

import (
	"fmt"
)

type Call struct {
	// XMLName     xml.Name `xml:"call"`
	PhoneNumber  string `xml:"number,attr"`
	Duration     int    `xml:"duration,attr"`
	Type         int    `xml:"type,attr"`
	Date         string `xml:"date,attr"`
	ReadableDate string `xml:"readable_date,attr"`
	ContactName  string `xml:"contact_name,attr"`
}

func (c Call) Format(f fmt.State, r rune) {
	fmt.Fprintf(f, "Call{%v from %v (%v)}", c.ReadableDate, c.ContactName, c.PhoneNumber)
}
