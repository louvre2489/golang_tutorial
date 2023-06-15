package main

import (
	"database/sql"
	"fmt"
	"log"
	//	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   "root",     // os.Getenv("DBUSER"),
		Passwd: "hogehoge", //os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

  album, err := albumByID(1)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Album found: %v\n", album)

  alb := Album {
    Title: "Test",
    Artist: "Sample",
    Price: 98.7,
  }
  id, err := addAlbum(alb)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("New Album ID: %d\n", id)
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func albumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %d: no such album.", id)
		}
		return alb, fmt.Errorf("albumByID %d: %v", id, err)
	}

  return alb, nil
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist =?", name)
	if err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}

	return albums, nil
}

func addAlbum(alb Album)(int64, error) {
  result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
  if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
  }

  id, err := result.LastInsertId()
  if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
  }

  return id, nil
}
