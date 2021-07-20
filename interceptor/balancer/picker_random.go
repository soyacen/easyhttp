package easyhttpbalancer

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

type IntnRand interface {
	Intn(n int) int
}

type RandomPicker struct {
	intnRand IntnRand
}

func (picker *RandomPicker) Pick(pickerInfo PickerInfo) (PickResult, error) {
	hosts := strings.Split(pickerInfo.URL.Host, ",")
	if len(hosts) <= 0 {
		return PickResult{}, errors.New("host list is empty")
	}
	index := picker.intnRand.Intn(len(hosts))
	return PickResult{Host: hosts[index]}, nil
}

func NewRandomPicker() *RandomPicker {
	return &RandomPicker{intnRand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
