package core

import (
	"database/sql"
	"fmt"
	"strings"
)

// InsertIntoDB tries to insert table into database
func InsertIntoDB(data Table, db *sql.DB) error {
	cols := strings.Join(GetStructFields(data), ", ")
	//TODO Escape insert
	insert := fmt.Sprintf(`insert ignore Ads
	(%s)
	values 
	(%s)
	`, cols, data)
	fmt.Println(insert)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(insert)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}
