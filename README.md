<style>
    .tob-wrapper{
        height: 30vh;
        overflow-y: 'scroll'
    }
</style>

<div class="tob-wrapper" markdown="1">

# Summary

[Getting started](#getting-started)

[Comparing JSON and CSV](#comparing-json-and-csv)

- [CSV](#csv)
- [JSON](#json)

[Things that I want to bring in TOC swagger](#things-that-i-want-to-bring-in-toc-swagger)

</div>

## Getting started

To run the application locally, please follow the steps:

1. start your docker engine or open your docker desktop
2. In the directory where `docker-compose.yaml` located, type in the command `docker compose up -d`.
3. the service will run on following ports
   - `localhost:3000` for swagger_doc
   - `localhost:8000` for goapi

> please ignore the mysqldb service and the endpoint under Users domain, those are my playground. The endpoint under user only work if you uncomment everything in `docker-compose.yaml` and `database.ConnectDB()` in main.go

![The swagger page under localhost:3000](rmImg/Screenshot%20from%202023-08-15%2010-08-51.png)

## Comparing JSON and CSV

### CSV

People have been using it around the lauch of Macintosh[^mac-history]. Every operating system provide a very good interface to edit this format of file. Enough background, let's dive into pros and cons for uploading devices using csv.

As mentioned in the previous paragraph, most of the operating system providing good interface, which the interface give us the sense of data entry and column name. In the example below, I intentionally type in an invalid longitude. The API then return which data entry has the error, this help user track where the error is. On the other hand, I can customize my error message in the code.

However, the downside is very obvious as well. Notice that I need to parse every single field in `string` type to corresponding type in the custom `struct` using `strconv` package. This make it harder to scale on developer side.

![Csv wrong data entry at row 3 column D](rmImg/Screenshot%20from%202023-08-15%2010-30-08.png)

![API return error message for csv endpoint](rmImg/Screenshot%20from%202023-08-15%2010-33-19.png)
[^mac-history]: Apple's computer after Lisa. After Steve being oust from Lisa team, he took over control of Apple's side project Macintosh which originally target not demanding user that use personal computer to do word processing and reading. It has huge impact on personal computer industry, much like Elon's Tesla vehicle to eletric car industry

```go
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

```

[<< back to top](#summary)

### JSON

Sorry I have not read about the man who create JSON, so not much background and history here. Let's dive right in.

If user like us edit with vscode, we did get an error when we type in invalid character for certain field. This is of course the first line of defense to invalid data type. The second line of defense is the `validator` package that was used a lot in our TOC API. The `validator` package give us an interface named `FieldError` that has serveral method extracting the error detail.

> I have not fully finish experimenting all method of `FieldError`, so there might still be method that we can use to customize our error message.

The downside is relating to `validator` package as well. From my experiment, the method which is supposed to return the error message return very general detail the error, as shown below. It is smart enough to detect the invalid character, but it doesn't tell the end user which data entry contain the invalid data.

![invalid character warning in JSON](rmImg/Screenshot%20from%202023-08-15%2010-50-12.png)

![bad request reponse for JSON api](rmImg/Screenshot%20from%202023-08-15%2011-01-28.png)

[<< back to top](#summary)

## Things that I want to bring in TOC swagger

- `tag`
  : the tag tells the API user which model they are interacting with. Taking this api for example, we have Users and Devices. In the yaml file, it is root level properties, `tags` is an array of object, and the object has property `name` and `description`. The API under the http actions also has an array property `tags`, we can use it to specify the defined tags which the operation belongs to.

  ```yaml
  tags:
    - name: Users
      description: users of the systems
    - name: Devices
      description: test bulk upload functionality of gNb device
  ```

  ```yaml
  /api/devices/csv:
    post:
      tags:
        - Devices
  ```

- `description`:
  This property has the same level as root level tags. Given the complexity of TOC, I think it is a good idea to tell our user how to interact with the API. Further more it support `markdown` syntax in yaml file. It will make our API page look nicer.

  ![markdown for swagger](rmImg/Screenshot%20from%202023-08-15%2011-19-18.png)

  ![swagger description](rmImg/Screenshot%20from%202023-08-15%2011-20-14.png)

  [<< back to top](#summary)
