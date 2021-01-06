package imoti

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// FetchXpaths gets url content and returns data from xpath locations
func FetchXpaths(url string, xpathJSONPaths ...string) (map[string]string, error) {
	node := contentToNode(url)
	if node == nil {
		return nil, errors.New("Invalid html")
	}
	data := make(map[string]string)
	for i := range xpathJSONPaths {
		jsonObj, err := deserializeXpathsJSON(xpathJSONPaths[i])
		var xpathMap map[string]string
		var required []string
		json.Unmarshal(jsonObj["XPaths"], &xpathMap)
		json.Unmarshal(jsonObj["Required"], &required)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for key, val := range xpathMap {
			data[key] = fetchByXpath(node, val)
		}

		passes := true
		for i := range required {
			if data[required[i]] == "" {
				passes = false
				break
			}
		}
		if passes {
			return data, nil
		}
	}
	return data, errors.New("Content didn't meet requirements")
}
func deserializeXpathsJSON(file string) (jsonObj map[string]json.RawMessage, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &jsonObj)
	return
}

func fetchByXpath(node *html.Node, xpath string) string {
	queried, err := htmlquery.Query(node, xpath)
	if err != nil {
		return ""
	}
	// if XPath is incorrect
	if queried == nil || queried.FirstChild == nil {
		return ""
	}
	// Get content of tag
	return queried.FirstChild.Data
}

func contentToNode(url string) (node *html.Node) {
	source, err := http.Get(url)
	if err != nil {
		return
	}
	// Check html received
	// data, err := ioutil.ReadAll(source.Body)
	// ioutil.WriteFile("obiava.html", data, os.FileMode(os.O_CREATE))

	node, err = html.Parse(source.Body)
	return
}
