package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"

	"./core"
	"./imoti"

	// MYSQL driver
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "crawler@/realestate")
	if err != nil {
		panic(err)
	}
	navpage, err := getNavpage()
	if err != nil {
		panic(err)
	}
	crawler := imoti.NewCrawler()
	go func() {
		err := crawler.FetchLinks(navpage)
		if err != nil {
			panic(err)
		}
	}()
	insertLinks(db, crawler.Links)
	db.Close()
}

// Reads local `crawl.json` for navpage link
func getNavpage() (string, error) {
	data, err := ioutil.ReadFile("./crawl.json")
	if err != nil {
		return "", err
	}
	jsonObj := struct {
		Navpage string
	}{}
	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		return "", err
	}
	return jsonObj.Navpage, err
}
func insertLinks(db *sql.DB, links <-chan string) {
	rand.Seed(time.Now().UnixNano())
	tableChan := make(chan core.Table)
	length := 0
	for {
		link, timedout := timeOutput(links, time.Second*7)
		if timedout {
			break
		}
		go func() {
			err := imoti.Parse(link, tableChan)
			if err != nil {
				fmt.Println(err)
			}
		}()
		length++
	}

	var wg sync.WaitGroup
	wg.Add(length)
	for i := 0; i < length; i++ {
		// foreach link try insert into database table
		go func() {
			data := <-tableChan
			err := core.InsertIntoDB(data, db)
			if err != nil {
				fmt.Println(err)
			}
			// fmt.Println(data)
			wg.Done()
		}()
	}
	wg.Wait()
	close(tableChan)
}

// times output of channel to duration
func timeOutput(links <-chan string, duration time.Duration) (result string, timedout bool) {
	var res chan string = make(chan string)
	timedout = false
	go func() {
		res <- <-links
	}()
	go func() {
		time.Sleep(duration)
		res <- ""
		timedout = true
	}()
	result = <-res
	return result, timedout
}
