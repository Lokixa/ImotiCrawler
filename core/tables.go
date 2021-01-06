package core

import (
	"fmt"
)

// Table is a data layout for real estate sites
type Table struct {
	Name      adName
	Location  location
	Price     string
	Activated string
	Agency    string
	Broker    string
	TypeOfAd  string
	URL       string
}

func (table Table) String() string {
	//TODO Refactor with reflection
	return fmt.Sprintf("%s, %s, %q, %q, %q, %q, %q, %q", table.Name, table.Location, table.Price, table.Activated, table.Agency, table.Broker, table.TypeOfAd, table.URL)
}

type location struct {
	Area   string
	Region string
}

func (loc location) String() string {
	return fmt.Sprintf("%q, %q", loc.Area, loc.Region)
}

type adName struct {
	Type string
	Size string
}

func (name adName) String() string {
	return fmt.Sprintf("%q, %q", name.Type, name.Size)
}
