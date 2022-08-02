package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := app.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM albums")
	app.DB.Exec("ALTER SEQUENCE albums_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS albums
(
	id SERIAL,
	artist TEXT NOT NULL,
	title TEXT NOT NULL,
	CONSTRAINT albums_pkey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/albums", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentAlbum(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/album/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Album not found" {
		t.Errorf(
			"Expected the 'error' key of the response body to be set to 'Album not found'. Got '%s'",
			m["error"],
		)
	}
}

func TestCreateAlbum(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"title":"test album","artist":"test artist"}`)
	req, _ := http.NewRequest("POST", "/album", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["artist"] != "test artist" {
		t.Errorf("Expected album artist to be 'test artist'. Got '%s'", m["artist"])
	}

	if m["title"] != "test album" {
		t.Errorf("Expected album title to be 'test album'. Got '%s'", m["title"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected album ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetAlbum(t *testing.T) {
	clearTable()
	addAlbums(1)

	req, _ := http.NewRequest("GET", "/album/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addAlbums(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		app.DB.Exec("INSERT INTO albums(artist,title) VALUES($1, $2)", "Artist "+strconv.Itoa(i), "Album "+strconv.Itoa(i))
	}
}
