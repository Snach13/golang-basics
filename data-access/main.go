package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
    ID int64
    Title string
    Artist string
    Price float32
}

func main() {

	// 	// Get environment variables
	// dbUser := os.Getenv("DBUSER")
	// dbPass := os.Getenv("DBPASS")

	// // Print the values to check if they are being read correctly
	// fmt.Println("DBUSER:", dbUser)
	// fmt.Println("DBPASS:", dbPass)

	// if dbUser == "" || dbPass == "" {
	// 	log.Fatal("Error: DBUSER or DBPASS is not set")
	// }

    // Capture connection properties.
    cfg := mysql.Config{
        User: "root",
        Passwd: "Nachi@123",
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "recordings",
        AllowNativePasswords: true,

    }
    // Get a database handle.
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
    fmt.Printf("Albums by John Coltrane: %v\n", albums)

    alb, err := albumById(2)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Album found: %v\n", alb)

    albID, err := addAlbum(Album{
        Title: "The Modern sound of betty carter",
        Artist: "Betty Carter",
        Price: 49.55,
    })

    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
    // An albums slice to hold data from returned rows.
    var albums []Album

    rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
    if err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    defer rows.Close()
    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var alb Album
        if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        albums = append(albums, alb)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    return albums, nil
}

func albumById(id int64) (Album, error) {
    var alb Album

    row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
    if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
        if err == sql.ErrNoRows {
            return alb, fmt.Errorf("albumById %d: %v", id, err)
        }

        return alb, fmt.Errorf("albumById %d: %v", id, err)
    }

    return alb, nil
}

func addAlbum(alb Album) (int64, error) {
    result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
    if err != nil {
        return 0, fmt.Errorf("add album: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("add album: %v", err)
    }

    return id, nil
}
