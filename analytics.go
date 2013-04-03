package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	BatchInsertSeconds int    `json:"batchInsertSeconds"`
}

func ReadAppConfig() AppConfig {
	appConfig := AppConfig{}
	appConfigFile, err := ioutil.ReadFile("config/app.json")
	if err != nil {
	        log.Fatal("Unable to read config file: ", err)
	}
	if err = json.Unmarshal(appConfigFile, &appConfig); err != nil {
	        log.Fatal("Unable to unmarshal appConfigFile into appConfig: ", err)
	}
	return appConfig
}

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

func hrefClickHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	hrefClick := HrefClick{}
	if err := json.Unmarshal(body, &hrefClick); err != nil {
		log.Println("Unable to unmarshal hrefClick: ", err)
	}
	// Get ip address from http request
	hrefClick.IpAddress = IpAddress(r.RemoteAddr)
	hrefClicks = append(hrefClicks, hrefClick)
	hrefClick.Status = "ok"
	responseJson, _ := json.Marshal(hrefClick)
	fmt.Fprintf(w, string(responseJson))
}

func pageViewsHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	pageView := PageView{}
	if err := json.Unmarshal(body, &pageView); err != nil {
		log.Println("Unable to unmarshal pageView: ", err)
	}
	// Get ip address from http request
	pageView.IpAddress = IpAddress(r.RemoteAddr)
	pageViews = append(pageViews, pageView)
	pageView.Status = "ok"
	responseJson, _ := json.Marshal(pageView)
	fmt.Fprintf(w, string(responseJson))
}

func ListenForRecords() {
	appConfig := ReadAppConfig()
	seconds := time.Duration(appConfig.BatchInsertSeconds)*time.Second
	// Run every x seconds.
	for _ = range time.Tick(seconds) {
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
			w.WriteHeader(405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Unable to read requeset body: ", err)
		}
		fn(w, r, body)
	}
}

func main() {
	logfile, err := os.OpenFile("var/error.log",
	        os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
	        log.Fatal("Unable to open log file: ", err)
	}
	log.SetOutput(logfile)
	http.HandleFunc("/page-views/", makeHandler(pageViewsHandler))
	http.HandleFunc("/href-click/", makeHandler(hrefClickHandler))
	go ListenForRecords()
	http.ListenAndServe(":8080", nil)
}
