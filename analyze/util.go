package analyze

import (
	// "fmt"
	"log"
	"sort"
	"time"
)

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

func GetSortedKeysInt64(m map[int64]float64) []int64 {
	keys := make([]int64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })
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

// Returns keys, values as two slices
func SplitMapSortedInt64(m map[int64]float64) ([]int64, []float64) {
	keys := GetSortedKeysInt64(m)
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

func CountKeyInt64(m map[int64]float64, key int64, amount float64) {
	if count, ok := m[key]; ok {
		m[key] = count + amount
	} else {
		m[key] = amount
	}
}
