package analyze

import (
	"log"
	"math"
)

type Histogram struct {
	data           []int64
	BinGranularity int
	MaxBins        int
}

func MakeHistogram() Histogram {
	return Histogram{
		data:           []int64{},
		BinGranularity: 80,
		MaxBins:        20,
	}
}

func (h *Histogram) Add(val int64) {
	h.data = append(h.data, val)
}

func (h *Histogram) GetMap() map[int64]float64 {
	bin_size := h.getBinSize()
	result := map[int64]float64{}
	for _, val := range h.data {
		bin_index := val / bin_size
		if bin_index >= int64(h.MaxBins-1) {
			bin_index = int64(h.MaxBins - 1)
		}
		bin := bin_index * bin_size
		CountKeyInt64(result, bin, 1)
	}
	return result
}

func (h *Histogram) getBinSize() int64 {
	var min int64 = math.MaxInt64
	var max int64 = math.MinInt64
	if len(h.data) == 0 {
		return 1
	}

	for _, val := range h.data {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}

	bin_size := (max - min) / int64(h.BinGranularity)
	log.Printf("max = %v", max)
	log.Printf("min = %v", min)
	log.Printf("bin_size = %v", bin_size)
	if bin_size == 0 {
		bin_size = 1
	} else if bin_size < 0 {
		bin_size = -1 * bin_size
	}
	return bin_size
}
