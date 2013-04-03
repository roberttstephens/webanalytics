package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

func SetPageViews(p []PageView) {
	if len(p) < 1 {
		return
	}
	db := Db()
	tx, _ := db.Begin()
	stmt, err := db.Prepare("INSERT INTO page_view(timestamp, url, ip_address, user_agent, screen_height, screen_width) VALUES (NOW(), $1, $2, $3, $4, $5)")
	if err != nil {
		log.Println("Unable to prepare statment for PageView: ", err)
	}
	for _, p := range p {
		tx.Stmt(stmt).Exec(p.Url, p.IpAddress, p.UserAgent, p.ScreenHeight, p.ScreenWidth)
	}
	tx.Commit()
}

func SetHrefClicks(h []HrefClick) {
	if len(h) < 1 {
		return
	}
	db := Db()
	tx, _ := db.Begin()
	stmt, err := db.Prepare("INSERT INTO href_click(timestamp, url, ip_address, href, href_rectangle) VALUES (NOW(), $1, $2, $3, box(point($4,$5), point($6,$7)))")
	if err != nil {
		log.Println("Unable to prepare statment for HrefClick: ", err)
	}
	for _, h := range h {
		tx.Stmt(stmt).Exec(h.Url, h.IpAddress, h.Href, h.HrefTop, h.HrefRight, h.HrefBottom, h.HrefLeft)
	}
	tx.Commit()
}
