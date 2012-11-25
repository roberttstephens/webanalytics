package db

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
)

// TODO maybe add a setPageView method and setHhrefClick method.
// The set setPageView method will read from app.yaml and will determine how
// many page views and href clicks to store in memory, and a maximum timeout

func main() {
	db, err := sql.Open("postgres", "user=analytics password=analytics host=127.0.0.1 dbname=analytics")
	if err != nil {
		fmt.Println("Error")
	}
	rows, err := db.Query("select * from page_view")
	if err != nil {
		fmt.Println("Error 2")
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
