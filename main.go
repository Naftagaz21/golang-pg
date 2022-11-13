package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 9000
	user     = "postgres"
	password = "atanasj"
	dbname   = "MyFavouriteColors"
)

type Color struct {
	id       int64
	title    string
	colorhex string
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println(db)
	fmt.Println(("Successfully connected!"))

	rows, err := db.Query(`SELECT * FROM colors`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var data []Color
	for rows.Next() {
		var (
			id       int64
			name     string
			colorHex string
		)
		//rows.Scan()
		rows.Scan(&id, &name, &colorHex)
		data = append(data, Color{id: id, title: name, colorhex: colorHex})
		//todos = append(todos, res)
	}

	fmt.Println((data))

}
