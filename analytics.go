package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PageView struct {
	IpAddress  string `json:"ipAddress"`
	ScreenSize string `json:"screenSize"`
	Timestamp  int    `json:"timestamp"`
	Url        string `json:"url"`
	UserAgent  string `json:"userAgent"`
}

func pageViewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	fmt.Printf("Post received\n")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Readall Error: ", err)
		return
	}
	pageView := PageView{}
	err = json.Unmarshal(body, &pageView)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Results: %v\n", pageView)
	fmt.Printf("Ip Address: %v\n", pageView.IpAddress)
	// TODO insert into db
}

func hrefClickHandler(w http.ResponseWriter, r *http.Request) {
	// TODO parse JSON from post, insert into DB.
}

func main() {
	http.HandleFunc("/page-views/", pageViewsHandler)
	http.HandleFunc("/href-click/", hrefClickHandler)
	ReadConfig()
	http.ListenAndServe(":8080", nil)
}
