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
	dbConfigFile, err := ioutil.ReadFile("db_config.json")
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

func setPageView() {
  db := Db()
  sql := fmt.Sprintf("INSERT INTO page_view(ip_address, url, time, browser) VALUES ('%s', '%s', NOW(), '%s')",
    

}
