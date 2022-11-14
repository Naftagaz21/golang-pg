package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	ColorHex string `json:"colorhex"`
}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func EnableCors(res *http.ResponseWriter) {
	//(*res).Header().Set("Access-Control-Allow-Origin", "*")
	//(*res).Header().Set("Content-Type", "text/html; charset=utf-8")
	//(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
}

func GETHandler(res http.ResponseWriter, req *http.Request) {
	EnableCors(&res)
	db := OpenConnection()

	rows, err := db.Query(`SELECT * FROM colors`)
	if err != nil {
		log.Fatal()
	}

	var colors []Color

	for rows.Next() {
		var color Color
		rows.Scan(&color.Id, &color.Title, &color.ColorHex)
		colors = append(colors, color)
	}

	colorsBytes, _ := json.Marshal(colors)
	res.Header().Set("Content-Type", "application/json")
	res.Write(colorsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(res http.ResponseWriter, req *http.Request) {
	EnableCors(&res)
	db := OpenConnection()

	var color Color
	err := json.NewDecoder(req.Body).Decode(&color)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	insertQuery := `INSERT INTO colors (title, hexcode) VALUES ($1, $2)`
	_, err = db.Exec(insertQuery, color.Title, color.ColorHex)
	// Returns proper http message based on whether duplicates are found
	if err != nil {
		if strings.HasPrefix(err.Error(), `pq: duplicate key`) {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusBadRequest)
			resp := make(map[string]string)
			if strings.HasSuffix(err.Error(), `"unique_title"`) {
				resp["message"] = "Duplicate color name"
			} else {
				resp["message"] = "Duplicate color code"
			}
			jsonResp, _ := json.Marshal(resp)
			res.Write(jsonResp)
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
		panic(err)
	}

	res.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {
	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
