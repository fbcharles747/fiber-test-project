package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/fbcharles747/fiber-api/models"
	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	ErrMsg string `json:"errMsg"`
}

func BulkUploadCSV(c *fiber.Ctx) error {
	// get form field "file"
	fileheader, err := c.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	fileptr, err := fileheader.Open()
	defer fileptr.Close()

	if err != nil {
		resErr := ErrorResponse{ErrMsg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(resErr)
	}

	var result *[]models.ResponseDevice

	if filepath.Ext(fileheader.Filename) == ".csv" {
		csvDevices, err := CsvParser(&fileptr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}
		result, err = CreateResponseDevices(csvDevices)
		if err != nil {
			resErr := ErrorResponse{ErrMsg: err.Error()}
			return c.Status(http.StatusBadRequest).JSON(resErr)
		}

	} else {
		responseErr := ErrorResponse{ErrMsg: "Wrong file format"}
		return c.Status(http.StatusBadRequest).JSON(responseErr)
	}

	return c.Status(http.StatusCreated).JSON(result)

}

func CsvParser(csvFile *multipart.File) ([]*models.CsvDevice, error) {
	csvDevices := []*models.CsvDevice{}
	if err := gocsv.UnmarshalMultipartFile(csvFile, &csvDevices); err != nil {
		return nil, errors.New("not able to open the csv file")
	}

	return csvDevices, nil
}

func CreateResponseDevices(devices []*models.CsvDevice) (result *[]models.ResponseDevice, err error) {
	var resultSlice []models.ResponseDevice
	for i := 0; i < len(devices); i++ {
		var responseDevice models.ResponseDevice
		responseDevice.DeviceId = devices[i].DeviceId
		responseDevice.DeviceType = devices[i].DeviceType
		responseDevice.Latitude, err = strconv.ParseFloat(devices[i].Latitude, 32)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot parse Data entry %d, original error message: %s", i+1, err.Error())
			return nil, errors.New(errMsg)
		}
		responseDevice.Longitude, err = strconv.ParseFloat(devices[i].Longitude, 32)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot parse Data entry %d, original error message: %s", i+1, err.Error())
			return nil, errors.New(errMsg)
		}
		responseDevice.StreetAddress = devices[i].StreetAddress
		resultSlice = append(resultSlice, responseDevice)
	}

	result = &resultSlice

	return result, nil

}

// below is for json file
// ====================================================================================

type ErrMap map[string][]models.IError

func BulkUploadJSON(c *fiber.Ctx) error {
	// parse the file
	fileheader, err := c.FormFile("file")
	if err != nil {
		errRes := ErrorResponse{ErrMsg: "File not found"}
		return c.Status(http.StatusBadRequest).JSON(errRes)
	}
	filePtr, err := fileheader.Open()
	if err != nil {
		errRes := ErrorResponse{ErrMsg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(errRes)
	}
	defer filePtr.Close()

	// decode the file content
	decoder := json.NewDecoder(filePtr)

	var deviceSlice []models.ResponseDevice

	if err := decoder.Decode(&deviceSlice); err != nil {
		errRes := ErrorResponse{ErrMsg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(errRes)
	}

	// validate the file
	errMap := validateResponseDevice(&deviceSlice)
	if len(errMap) == 0 {
		return c.Status(http.StatusCreated).JSON(deviceSlice)
	} else {
		return c.Status(http.StatusBadRequest).JSON(errMap)
	}

}

func validateResponseDevice(responseDevices *[]models.ResponseDevice) ErrMap {
	errMap := make(ErrMap)

	for _, device := range *responseDevices {
		err := device.Validate()
		if errs, ok := err.(validator.ValidationErrors); ok {
			// errs is now a []FieldError, which contain all the validation error belonging to one data entry
			var fieldErrs []models.IError

			for _, fieldErr := range errs {
				var errItem models.IError
				errItem.Error = fieldErr.Error()
				errItem.Field = fieldErr.Field()
				fieldErrs = append(fieldErrs, errItem)
			}
			errMap[device.DeviceId] = fieldErrs
		}
	}
	return errMap
}
