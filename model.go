package main

import (
	"database/sql"
)

// Representation of Album table in database
type album struct {
	ID     int    `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

// Query database for Album by ID
func (a *album) getAlbum(db *sql.DB) error {
	return db.QueryRow("SELECT artist, title FROM albums WHERE id=$1", a.ID).Scan(&a.Artist, &a.Title)
}

// Updates Album in database with new artist and title
func (a *album) updateAlbum(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE albums SET artist=$1, title=$2 WHERE id=$3",
			a.Artist, a.Title, a.ID)
	return err
}

// Deletes an Album in database by ID
func (a *album) deleteAlbum(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM albums WHERE id=$1", a.ID)

	return err
}

// Inserts Album into database
func (a *album) createAlbum(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO albums(artist, title) VALUES($1, $2) RETURNING id",
		a.Artist, a.Title).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}

// Select all from Albums in database
func getAlbums(db *sql.DB, start, count int) ([]album, error) {
	rows, err := db.Query(
		"SELECT id, artist, title FROM albums LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	albums := []album{}

	for rows.Next() {
		var a album
		if err := rows.Scan(&a.ID, &a.Artist, &a.Title); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, nil
}
