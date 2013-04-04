package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/bmizerany/pq"
	"io/ioutil"
	"log"
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
		log.Fatal("Unable to read config/db.json: ", err)
	}
	if err = json.Unmarshal(dbConfigFile, &dbConfig); err != nil {
		log.Fatal("Unable to unmarshal dbConfig: ", err)
	}
	return dbConfig
}

func Db() *sql.DB {
	dbConfig := ReadDbConfig()
	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s host=%s dbname=%s", dbConfig.User,
			dbConfig.Pass, dbConfig.Host, dbConfig.Name))
	if err != nil {
		log.Fatal("Unable to connect to the database: ", err)
	}
	return db
}

func PageViews() {
	db := Db()
	rows, err := db.Query("select * from page_view")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, time int
		var ip_address, url, browser string
		rows.Scan(&id, &ip_address, &url, &time, &browser)
		fmt.Println(url)
	}
}

func SetPageViews(pageViews []PageView) {
	if len(pageViews) < 1 {
		return
	}
	db := Db()
	tx, _ := db.Begin()
	stmt, err := db.Prepare("INSERT INTO page_view(timestamp, url, ip_address, user_agent, screen_height, screen_width) VALUES (NOW(), $1, $2, $3, $4, $5)")
	if err != nil {
		log.Println("Unable to prepare statment for PageView: ", err)
	}
	for k := range pageViews {
		tx.Stmt(stmt).Exec(
			pageViews[k].Url,
			pageViews[k].IpAddress,
			pageViews[k].UserAgent,
			pageViews[k].ScreenHeight,
			pageViews[k].ScreenWidth,
		)
	}
	tx.Commit()
}

func SetHrefClicks(hrefClicks []HrefClick) {
	if len(hrefClicks) < 1 {
		return
	}
	db := Db()
	tx, _ := db.Begin()
	stmt, err := db.Prepare("INSERT INTO href_click(timestamp, url, ip_address, href, href_rectangle) VALUES (NOW(), $1, $2, $3, box(point($4,$5), point($6,$7)))")
	if err != nil {
		log.Println("Unable to prepare statment for HrefClick: ", err)
	}
	for k := range hrefClicks {
		tx.Stmt(stmt).Exec(
			hrefClicks[k].Url,
			hrefClicks[k].IpAddress,
			hrefClicks[k].Href,
			hrefClicks[k].HrefTop,
			hrefClicks[k].HrefRight,
			hrefClicks[k].HrefBottom,
			hrefClicks[k].HrefLeft,
		)
	}
	tx.Commit()
}
