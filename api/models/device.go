package models

import (
	"github.com/go-playground/validator/v10"
)

type CsvDevice struct {
	DeviceId      string `csv:"device_id"`
	DeviceType    string `csv:"device_type"`
	Latitude      string `csv:"latitude"`
	Longitude     string `csv:"longitude"`
	StreetAddress string `csv:"street_address"`
}

type ResponseDevice struct {
	DeviceId      string  `json:"device_id";`
	DeviceType    string  `json:"device_type" `
	Latitude      float64 `json:"latitude" validate:"min=-8850,max=8850"`
	Longitude     float64 `json:"longitude" validate:"min=-8850,max=8850"`
	StreetAddress string  `json:"street_address"`
}

type IError struct {
	Field string `json:"field"`
	Error string `json:"errMsg"`
}

func (d *ResponseDevice) Validate() error {
	validator := validator.New()
	return validator.Struct(d)
}
