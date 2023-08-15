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

func Check(err error) {
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}

type ErrMap map[string][]models.IError

func badrequestHandler(c *fiber.Ctx) error {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			errRes := ErrorResponse{ErrMsg: err.Error()}
			return c.Status(http.StatusBadRequest).JSON(errRes)
		}
	}
	errRes := ErrorResponse{ErrMsg: "unknown"}
	return c.Status(http.StatusBadRequest).JSON(errRes)
}

func BulkUploadJSON(c *fiber.Ctx) error {

	// parse the file
	fileheader, err := c.FormFile("file")
	if err != nil {
		resErr := ErrorResponse{ErrMsg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(resErr)
	}
	if filepath.Ext(fileheader.Filename) != ".json" {
		responseErr := ErrorResponse{ErrMsg: "Wrong file format, this end point only accept JSON"}
		return c.Status(http.StatusBadRequest).JSON(responseErr)
	}
	filePtr, err := fileheader.Open()
	if err != nil {
		resErr := ErrorResponse{ErrMsg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(resErr)
	}
	defer filePtr.Close()

	// decode the file content
	decoder := json.NewDecoder(filePtr)
	var deviceSlice []models.ResponseDevice
	// need to reset err here
	err = nil
	err = decoder.Decode(&deviceSlice)
	if err != nil {
		resErr := ErrorResponse{ErrMsg: err.Error()}
		return c.Status(http.StatusBadRequest).JSON(resErr)
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
