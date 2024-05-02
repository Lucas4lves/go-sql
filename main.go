package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Album struct {
	Title  string
	Artist string
	Price  float64
}

func CreateTable(db *sql.DB) {
	sql := `CREATE TABLE IF NOT EXISTS albums
				(id SERIAL PRIMARY KEY,
				 title VARCHAR(50) NOT NULL,
				 artist VARCHAR(50) NOT NULL,
				 price DECIMAL(10,2) NULL,
				 created timestamp DEFAULT NOW()
				 )`

	_, err := db.Exec(sql)

	if err != nil {
		log.Fatalf("ERR CreateTable: %s", err)
	}
}

func InsertAlbum(db *sql.DB, album Album) int {
	sql := `INSERT INTO albums (title, artist, price)
				VALUES($1, $2, $3)
				RETURNING id`

	var primaryKey int
	err := db.QueryRow(sql, album.Title, album.Artist, album.Price).Scan(&primaryKey)
	if err != nil {
		log.Fatalf("ERR: %s", err)
	}
	fmt.Printf("BEFORE : %d\n", primaryKey)
	return primaryKey
}

func ExecStatement(db *sql.DB, stmt string) {
	_, err := db.Exec(stmt)

	if err != nil {
		log.Fatalf("ERR: %s", err)
	}
}

func GetAlbumById(db *sql.DB, id int) Album {

	sql := "SELECT title, artist, price FROM albums WHERE id=$1"

	var title string
	var artist string
	var price float64

	err := db.QueryRow(sql, id).Scan(&title, &artist, &price)

	if err != nil {
		log.Fatalf("ERR: %s", err)
	}

	return Album{
		Title:  title,
		Artist: artist,
		Price:  price,
	}
}

func GetManyAlbums(db *sql.DB) []Album {
	data := []Album{}

	rows, err := db.Query("SELECT title, artist, price FROM albums")

	if err != nil {
		log.Fatal(err)
	}

	var title string
	var artist string
	var price float64

	for rows.Next() {
		err := rows.Scan(&title, &artist, &price)

		if err != nil {
			log.Fatal(err)
		}

		o := Album{Title: title, Artist: artist, Price: price}

		data = append(data, o)
	}

	return data

}

func main() {
	connStr := "postgres://postgres:1234@localhost:5432/gotest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	res := GetManyAlbums(db)

	for _, a := range res {
		fmt.Println(a.Title, a.Artist, a.Price)
	}

	defer db.Close()
}
