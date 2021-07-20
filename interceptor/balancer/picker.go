package easyhttpbalancer

import (
	"context"
	"net/http"
	"net/url"
)

type PickerInfo struct {
	URL    *url.URL
	Header http.Header
	Ctx    context.Context
}

type PickResult struct {
	Host string
}

type Picker interface {
	Pick(pickerInfo PickerInfo) (PickResult, error)
}
