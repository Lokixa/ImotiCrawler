package imoti

import (
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"../core"
)

// Parse searches the node for xpaths defined in the imoti package
func Parse(url string, tableChan chan<- core.Table) (err error) {
	filters := map[string]func(string, *core.Table){
		"Name":      setName,
		"Location":  setLocation,
		"Activated": setActivated,
		"Price":     setPrice,
	}
	files, err := ioutil.ReadDir("imoti/adTypes/")
	if err != nil {
		return
	}
	var adTypes []string
	for i := range files {
		adTypes = append(adTypes, "imoti/adTypes/"+files[i].Name())
	}
	// Fetch all xpaths in json files
	xpaths, err := FetchXpaths(url, adTypes...)
	// if err != nil {
	// 	return
	// }

	table := core.Table{}
	table.URL = url

	tableVal := reflect.ValueOf(&table).Elem()
	tableSig := tableVal.Type()

	for i := 0; i < tableSig.NumField(); i++ {

		fieldName := tableSig.Field(i).Name
		// Actual field
		field := tableVal.FieldByName(fieldName)

		if xpaths[fieldName] != "" {
			if filters[fieldName] != nil {
				// Filters assign directly to table. Could refactor in future.
				filters[fieldName](xpaths[fieldName], &table)
			} else {
				// Assign to field
				field.SetString(xpaths[fieldName])
			}
		}
	}

	tableChan <- table
	return
}
func setPrice(price string, table *core.Table) {
	table.Price = strings.ReplaceAll(price, " ", "")
}
func setName(name string, data *core.Table) {
	if name == "" {
		return
	}
	adType := regexp.MustCompile(`.*,`)
	size := regexp.MustCompile(`, [0-9]*`)
	data.Name.Type = strings.Replace(adType.FindString(name), ",", "", 1)
	data.Name.Size = size.FindString(name)[2:]
}
func setLocation(location string, data *core.Table) {
	if location == "" {
		return
	} else if !strings.Contains(location, ",") {
		data.Location.Area = location
		return
	}

	city := regexp.MustCompile(`.*,`)
	street := regexp.MustCompile(`,.*`)
	data.Location.Area = strings.Replace(city.FindString(location), ",", "", 1)
	data.Location.Region = street.FindString(location)[2:]
}
func setActivated(activated string, data *core.Table) {
	if activated == "" {
		return
	}
	invExp := regexp.MustCompile(`[0-9]{2}.*`)
	data.Activated = invExp.FindString(activated)
}
