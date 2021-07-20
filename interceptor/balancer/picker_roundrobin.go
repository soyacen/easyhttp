package easyhttpbalancer

import (
	"errors"
	"strings"
	"sync"
)

type RoundRobinPicker struct {
	mu    sync.Mutex
	nexts map[string]int
}

func NewRoundRobinPicker() *RoundRobinPicker {
	return &RoundRobinPicker{nexts: make(map[string]int)}
}

func (picker *RoundRobinPicker) Pick(pickerInfo PickerInfo) (PickResult, error) {
	picker.mu.Lock()
	defer picker.mu.Unlock()
	hosts := strings.Split(pickerInfo.URL.Host, ",")
	if len(hosts) <= 0 {
		return PickResult{}, errors.New("host list is empty")
	}
	index := picker.nexts[pickerInfo.URL.Host]
	picker.nexts[pickerInfo.URL.Host] = (index + 1) % len(hosts)
	return PickResult{Host: hosts[index]}, nil
}
