package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
  "os"
  "time"
)

type PageView struct {
	IpAddress  string `json:"ipAddress"`
	ScreenSize string `json:"screenSize"`
	Timestamp  int    `json:"timestamp"`
	Url        string `json:"url"`
	UserAgent  string `json:"userAgent"`
}

func logError(s error) {
  f, err := os.OpenFile("var/error.log", os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  _, err = f.WriteString(fmt.Sprintf("%v\t%s\n", time.Now(), s))
  if err != nil {
    panic(err)
  }
}

func pageViewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError(err)
		return
	}
	pageView := PageView{}
	err = json.Unmarshal(body, &pageView)
	if err != nil {
		logError(err)
	}
  SetPageView(pageView)
}

func hrefClickHandler(w http.ResponseWriter, r *http.Request) {
	// TODO parse JSON from post, insert into DB.
}

func main() {
	http.HandleFunc("/page-views/", pageViewsHandler)
	http.HandleFunc("/href-click/", hrefClickHandler)
	http.ListenAndServe(":8080", nil)
}
