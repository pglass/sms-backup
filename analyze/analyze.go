package analyze

import (
	// "fmt"
	"log"
	"sort"
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

	a.countMessages()
}

type Message interface {
	GetTime() time.Time
	IsIncoming() bool
}

func (a *Analyzer) countMessages() {
	// Count the SMS messages
	for _, sms := range a.doc.SMSes {
		a.countMessage(sms)
	}
	for _, mms := range a.doc.MMSes {
		a.countMessage(mms)
	}
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

func PrintSorted(m map[time.Time]float64) {
	keys := GetSortedKeys(m)
	for _, k := range keys {
		log.Printf("%v = %v", k, m[k])
	}
}

func GetSortedKeys(m map[time.Time]float64) []time.Time {
	keys := make([]time.Time, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i].Unix() < keys[j].Unix() })
	return keys
}

// Returns keys, values as two slices
func SplitMapSorted(m map[time.Time]float64) ([]time.Time, []float64) {
	keys := GetSortedKeys(m)
	values := make([]float64, len(m))
	for i, k := range keys {
		values[i] = m[k]
	}
	return keys, values
}

func CountKey(m map[time.Time]float64, key time.Time, amount float64) {
	if count, ok := m[key]; ok {
		m[key] = count + amount
	} else {
		m[key] = amount
	}
}
