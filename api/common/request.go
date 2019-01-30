package common

import (
	"errors"
)

var (
	errMissingFilterKey  = errors.New("missing filter key")
	errInvalidFilterType = errors.New("invalid filter type")
)

type AssetProgram struct {
	Asset   string `json:"asset"`
	Program string `json:"program"`
}

// Display define how the data is displayed
type Display struct {
	Filter map[string]interface{} `json:"filter"`
	Sorter Sorter                 `json:"sort"`
}

type Sorter struct {
	By    string `json:"by"`
	Order string `json:"order"`
}

// GetFilterString give the filter keyword return the string value
func (d *Display) GetFilterString(filterKey string) (string, error) {
	if _, ok := d.Filter[filterKey]; !ok {
		return "", errMissingFilterKey
	}
	switch val := d.Filter[filterKey].(type) {
	case string:
		return val, nil
	}
	return "", errInvalidFilterType
}

// GetFilterInt give the filter keyword return the integer value
func (d *Display) GetFilterInt(filterKey string) (int, error) {
	if _, ok := d.Filter[filterKey]; !ok {
		return 0, errMissingFilterKey
	}
	switch val := d.Filter[filterKey].(type) {
	case int:
		return val, nil
	}
	return 0, errInvalidFilterType
}
