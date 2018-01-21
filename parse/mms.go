package parse

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// From the Android source, there's a PDUHeaders.java which contains:
//     public static final int FROM                            = 0x89;
//     public static final int TO                              = 0x97;
//
// https://android.googlesource.com/platform/frameworks/opt/mms/+/4bfcd8501f09763c10255442c2b48fad0c796baa/src/java/com/google/android/mms/pdu/PduHeaders.java
const (
	MMS_FROM_TYPE = 0x89
	MMS_TO_TYPE   = 0x97
)

type MMS struct {
	Parts []MMSPart `xml:"parts>part"`
	Addrs []MMSAddr `xml:"addrs>addr"`

	Date         int    `xml:"date,attr"`
	PhoneNumber  string `xml:"address,attr"`
	ReadableDate string `xml:"readable_date,attr"`
	ContactName  string `xml:"contact_name,attr"`

	my_number string
}

type MMSPart struct {
	ContentType string `xml:"ct,attr"`
	// This will contain the image data
	Data string `xml:"data,attr"`
	// This will contain text contents if there was a text message
	Text string `xml:"text,attr"`
}

type MMSAddr struct {
	Address string `xml:"address,attr"`
	// This indicates the sender and recipient
	Type int `xml:"type,attr"`
}

func (mms *MMS) SetMyNumber(number string) {
	mms.my_number = number
}

func (mms MMS) GetTime() time.Time {
	// sms.Date is a Unix timestamp in milliseconds
	seconds := int64(mms.Date / 1000)
	nanoseconds := int64((mms.Date % 1000) * 1000000)
	return time.Unix(seconds, nanoseconds).In(CENTRAL_LOCATION)
}

func (mms MMS) IsIncoming() bool {
	if mms.my_number == "" {
		log.Fatalf("Cannot determine incoming-ness for MMS (my_number = %v)", mms.my_number)
	}

	for _, addr := range mms.Addrs {
		if addr.Type == MMS_TO_TYPE && strings.Contains(addr.Address, mms.my_number) {
			return true
		} else if addr.Type == MMS_FROM_TYPE && strings.Contains(addr.Address, mms.my_number) {
			return false
		}
	}

	log.Printf("WARNNING: Could determine incoming-ness for MMS (addrs = %+v)", mms.Addrs)
	return false
}

func (mms MMS) Format(f fmt.State, r rune) {
	var from_or_to string
	if mms.IsIncoming() {
		from_or_to = "from"
	} else {
		from_or_to = "to"
	}

	var data_and_text string
	for _, part := range mms.Parts {
		if strings.Contains(part.ContentType, "image") {
			data_and_text += fmt.Sprintf("Image[<%v bytes>]", len(part.Data))
		} else if strings.Contains(part.ContentType, "text") {
			data_and_text += fmt.Sprintf("Text[%v]", part.Text)
		}
	}

	fmt.Fprintf(f, "MMS - %v %v %v: %v", mms.ReadableDate, from_or_to, mms.ContactName, data_and_text)
}
