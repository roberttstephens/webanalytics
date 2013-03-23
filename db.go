package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/bmizerany/pq"
	"io/ioutil"
)

type DbConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

func ReadDbConfig() DbConfig {
	dbConfig := DbConfig{}
	dbConfigFile, err := ioutil.ReadFile("config/db.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dbConfigFile, &dbConfig)
	if err != nil {
		panic(err)
	}
	return dbConfig
}

func Db() *sql.DB {
	dbConfig := ReadDbConfig()
	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s host=%s dbname=%s", dbConfig.User,
			dbConfig.Pass, dbConfig.Host, dbConfig.Name))
	if err != nil {
		panic(err)
	}
	return db
}

// TODO maybe add a setPageView method and setHhrefClick method.
// The set setPageView method will read from app.yaml and will determine how
// many page views and href clicks to store in memory, and a maximum timeout

func PageViews() {
	db := Db()
	rows, err := db.Query("select * from page_view")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var ip_address string
		var url string
		var time int
		var browser string
		rows.Scan(&id, &ip_address, &url, &time, &browser)
		fmt.Println(url)
	}
	rows.Close()
}

func SetPageViews() {
	if len(pageViews) < 1 {
		return
	}
	db := Db()
	tx, _ := db.Begin()
	stmt, err := db.Prepare("INSERT INTO page_view(timestamp, url, ip_address, user_agent, screen_height, screen_width) VALUES (NOW(), $1, $2, $3, $4, $5)")
	if err != nil {
		logError(err)
	}
	for _, p := range pageViews {
		tx.Stmt(stmt).Exec(p.Url, p.IpAddress, p.UserAgent, p.ScreenHeight, p.ScreenWidth)
	}
	// Empty the slice.
	pageViews = pageViews[0:0]
	tx.Commit()
}

func SetHrefClick(h HrefClick) {
	db := Db()
	stmt, err := db.Prepare("INSERT INTO href_click(timestamp, url, ip_address, href, href_rectangle) VALUES (NOW(), $1, $2, $3, box(point($4,$5), point($6,$7)))")
	if err != nil {
		logError(err)
	}
	_, err = stmt.Exec(h.Url, h.IpAddress, h.Href, h.HrefTop, h.HrefRight, h.HrefBottom, h.HrefLeft)
	if err != nil {
		logError(err)
	}
	fmt.Print("%v", h)
}
