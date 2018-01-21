package parse

import (
	"fmt"
	"log"
	"time"
)

type SMS struct {
	Address      string `xml:"address,attr"`
	Date         int    `xml:"date,attr"`
	Type         int    `xml:"type,attr"`
	Body         string `xml:"body,attr"`
	ReadableDate string `xml:"readable_date,attr"`
	ContactName  string `xml:"contact_name,attr"`
}

func (sms SMS) IsIncoming() bool {
	if sms.Type == 1 {
		return true
	} else if sms.Type == 2 {
		return false
	}
	log.Fatal("Unknown SMS Type: %v, body=%v, contact=%v", sms.Type, sms.Body, sms.ContactName)
	return false
}

func (sms SMS) Format(f fmt.State, r rune) {
	var from_or_to string
	if sms.IsIncoming() {
		from_or_to = "from"
	} else {
		from_or_to = "to"
	}
	fmt.Fprintf(f, "SMS - %v %s %v: %v", sms.ReadableDate, from_or_to, sms.ContactName, sms.Body)
}

func (sms SMS) GetTime() time.Time {
	// sms.Date is a Unix timestamp in milliseconds
	seconds := int64(sms.Date / 1000)
	nanoseconds := int64((sms.Date % 1000) * 1000000)
	return time.Unix(seconds, nanoseconds).In(CENTRAL_LOCATION)
}
