package analyze

import (
	// "fmt"
	"log"
	"time"

	"github.com/pglass/sms-backup/parse"
)

type Analyzer struct {
	doc parse.Document

	MessagesPerDay         map[time.Time]float64
	IncomingMessagesPerDay map[time.Time]float64
	OutgoingMessagesPerDay map[time.Time]float64

	MessagesPerWeek         map[time.Time]float64
	IncomingMessagesPerWeek map[time.Time]float64
	OutgoingMessagesPerWeek map[time.Time]float64

	IncomingMessageLengths Histogram
	OutgoingMessageLengths Histogram

	MessagesTimeOfDay map[time.Time]float64
}

func MakeAnalyzer(doc parse.Document) *Analyzer {
	return &Analyzer{doc: doc}
}

func (a *Analyzer) Run() {
	a.MessagesPerDay = map[time.Time]float64{}
	a.IncomingMessagesPerDay = map[time.Time]float64{}
	a.OutgoingMessagesPerDay = map[time.Time]float64{}

	a.MessagesPerWeek = map[time.Time]float64{}
	a.IncomingMessagesPerWeek = map[time.Time]float64{}
	a.OutgoingMessagesPerWeek = map[time.Time]float64{}

	a.IncomingMessageLengths = MakeHistogram()
	a.OutgoingMessageLengths = MakeHistogram()

	a.MessagesTimeOfDay = map[time.Time]float64{}

	// a.countMessages()
	// a.countMessageLengths()
	a.storeMessagesTimeOfDay()
}

type Message interface {
	GetTime() time.Time
	IsIncoming() bool

	GetContentType() int
	GetText() string
}

func (a *Analyzer) countMessages() {
	for _, sms := range a.doc.SMSes {
		a.countMessage(sms)
	}
	for _, mms := range a.doc.MMSes {
		a.countMessage(mms)
	}
}

func (a *Analyzer) countMessageLengths() {
	for _, sms := range a.doc.SMSes {
		a.countMessageLength(sms)
	}
	for _, mms := range a.doc.MMSes {
		a.countMessageLength(mms)
	}
}

func (a *Analyzer) storeMessagesTimeOfDay() {
	for _, sms := range a.doc.SMSes {
		a.storeMessageTimeOfDay(sms)
	}
	for _, mms := range a.doc.MMSes {
		a.storeMessageTimeOfDay(mms)
	}
}

func (a *Analyzer) storeMessageTimeOfDay(msg Message) {
	t := msg.GetTime()

	hour_as_fraction := float64(t.Hour()) + float64(t.Minute())/60.0

	a.MessagesTimeOfDay[t] = hour_as_fraction
}

func (a *Analyzer) countMessage(msg Message) {
	log.Printf("Processing: %v", msg)
	t := msg.GetTime()
	day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	week := time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
	if week.Weekday() != time.Sunday {
		log.Fatalf("Failed to floor datetime to current week")
	}

	// Count the total messages
	CountKey(a.MessagesPerDay, day, 1)
	CountKey(a.MessagesPerWeek, week, 1)

	// Count incoming and outgoing messages
	if msg.IsIncoming() {
		CountKey(a.IncomingMessagesPerDay, day, 1)
		CountKey(a.IncomingMessagesPerWeek, week, 1)
	} else {
		CountKey(a.OutgoingMessagesPerDay, day, 1)
		CountKey(a.OutgoingMessagesPerWeek, week, 1)
	}
}

func (a *Analyzer) countMessageLength(msg Message) {
	switch msg.GetContentType() {
	case parse.CONTENT_TEXT_AND_IMAGE, parse.CONTENT_TEXT:
		if msg.IsIncoming() {
			a.IncomingMessageLengths.Add(int64(len(msg.GetText())))
		} else {
			a.OutgoingMessageLengths.Add(int64(len(msg.GetText())))
		}
	}
}
