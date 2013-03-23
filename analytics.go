package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type HrefClick struct {
	IpAddress  string `json:"ipAddress"`
	Url        string `json:"url"`
	Href       string `json:"href"`
	HrefTop    int    `json:"hrefTop"`
	HrefRight  int    `json:"hrefRight"`
	HrefBottom int    `json:"hrefBottom"`
	HrefLeft   int    `json:"hrefLeft"`
	Status     string `json:"status"`
}

var hrefClicks []HrefClick

type PageView struct {
	Domain       string `json:"domain"`
	IpAddress    string `json:"ipAddress"`
	Url          string `json:"url"`
	UserAgent    string `json:"userAgent"`
	ScreenHeight int    `json:"screenHeight"`
	ScreenWidth  int    `json:"screenWidth"`
	Status       string `json:"status"`
}

var pageViews []PageView

func IpAddress(remoteAddr string) string {
	arr := strings.Split(remoteAddr, ":")
	return arr[0]
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

func hrefClickHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	hrefClick := HrefClick{}
	err := json.Unmarshal(body, &hrefClick)
	if err != nil {
		logError(err)
	}
	// Get ip address from http request
	hrefClick.IpAddress = IpAddress(r.RemoteAddr)
	hrefClicks = append(hrefClicks, hrefClick)
	hrefClick.Status = "ok"
	responseJson, err := json.Marshal(hrefClick)
	fmt.Fprintf(w, string(responseJson))
}

func pageViewsHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	pageView := PageView{}
	err := json.Unmarshal(body, &pageView)
	if err != nil {
		logError(err)
	}
	// Get ip address from http request
	pageView.IpAddress = IpAddress(r.RemoteAddr)
	pageViews = append(pageViews, pageView)
	pageView.Status = "ok"
	responseJson, err := json.Marshal(pageView)
	fmt.Fprintf(w, string(responseJson))
}

func SetRecords(d time.Duration) {
	// Run every d seconds.
	for _ = range time.Tick(d) {
		// Handle page views.
		newPageViews := make([]PageView, len(pageViews))
		copy(newPageViews, pageViews)
		go SetPageViews(newPageViews)
		pageViews = pageViews[0:0]

		// Handle href clicks.
		newHrefClicks := make([]HrefClick, len(hrefClicks))
		copy(newHrefClicks, hrefClicks)
		go SetHrefClicks(newHrefClicks)
		hrefClicks = hrefClicks[0:0]
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, []byte)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with, x-requested-by, Content-Type")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method != "POST" {
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logError(err)
			return
		}
		fn(w, r, body)
	}
}

func main() {
	http.HandleFunc("/page-views/", makeHandler(pageViewsHandler))
	http.HandleFunc("/href-click/", makeHandler(hrefClickHandler))
	go SetRecords(8 * time.Second)
	http.ListenAndServe(":8080", nil)
}
