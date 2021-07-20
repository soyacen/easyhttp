package easyhttpbalancer

import (
	"errors"
	"strings"
)

type FirstPicker struct{}

func (picker *FirstPicker) Pick(pickerInfo PickerInfo) (PickResult, error) {
	hosts := strings.Split(pickerInfo.URL.Host, ",")
	if len(hosts) <= 0 {
		return PickResult{}, errors.New("host list is empty")
	}
	return PickResult{Host: hosts[0]}, nil
}

func NewFirstPicker() *FirstPicker {
	return &FirstPicker{}
}
