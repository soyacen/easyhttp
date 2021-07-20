package easyhttpbalancer

import (
	"errors"
	"strings"

	"github.com/cespare/xxhash"
)

type KeyLocation bool

const (
	Header KeyLocation = true
	Query  KeyLocation = false
)

type HashPicker struct {
	key         string
	keyLocation KeyLocation
}

func (picker *HashPicker) Pick(pickerInfo PickerInfo) (PickResult, error) {
	hosts := strings.Split(pickerInfo.URL.Host, ",")
	if len(hosts) <= 0 {
		return PickResult{}, errors.New("host list is empty")
	}
	var val string
	if picker.keyLocation == Header {
		val = pickerInfo.Header.Get(picker.key)
	} else {
		val = pickerInfo.URL.Query().Get(picker.key)
	}
	sum := xxhash.Sum64String(val)
	index := int(sum) % len(hosts)
	return PickResult{Host: hosts[index]}, nil
}

func NewHashPicker(key string, keyLocation KeyLocation) *HashPicker {
	return &HashPicker{
		key:         key,
		keyLocation: keyLocation,
	}
}
