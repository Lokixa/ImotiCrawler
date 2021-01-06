package imoti

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"../core"
	"./xpaths"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// Parse searches the node for xpaths defined in the imoti package
func Parse(url string, tableChan chan core.Table) error {
	node := contentToNode(url)
	if node == nil {
		fmt.Println("Null html node")
		return errors.New("Invalid html")
	}
	var table core.Table

	table.URL = url
	setName(node, &table)

	setLocation(node, &table)

	price := fetchXpath(node, xpaths.Price)
	table.Price = strings.ReplaceAll(price, " ", "")

	setActivated(node, &table)

	table.Agency = fetchXpath(node, xpaths.Agency)

	table.Broker = fetchXpath(node, xpaths.Broker)

	table.TypeOfAd = fetchXpath(node, xpaths.TypeOfAd)

	tableChan <- table
	return nil
}
func setName(node *html.Node, data *core.Table) {
	name := fetchXpath(node, xpaths.Name)
	if name == "" {
		return
	}
	adType := regexp.MustCompile(`.*,`)
	size := regexp.MustCompile(`, [0-9]*`)
	data.Name.Type = strings.Replace(adType.FindString(name), ",", "", 1)
	data.Name.Size = size.FindString(name)[2:]
}
func setLocation(node *html.Node, data *core.Table) {
	location := fetchXpath(node, xpaths.Location)
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
func setActivated(node *html.Node, data *core.Table) {
	activated := fetchXpath(node, xpaths.Activated)
	if activated == "" {
		return
	}
	invExp := regexp.MustCompile(`[0-9]{2}.*`)
	data.Activated = invExp.FindString(activated)
}

func fetchXpath(node *html.Node, xpath string) string {
	queried, _ := htmlquery.Query(node, xpath)
	// if XPath is correct
	if queried != nil && queried.FirstChild != nil {
		// Get content of tag
		return queried.FirstChild.Data
	}
	return ""
}

func contentToNode(url string) (node *html.Node) {
	source, err := http.Get(url)
	if err != nil {
		return
	}
	node, err = html.Parse(source.Body)
	return
}
