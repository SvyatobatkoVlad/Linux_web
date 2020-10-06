package main

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSLLMode)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM book")
	a.DB.Exec("ALTER SEQUENCE book_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS book
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT book_pkey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/book", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentBook(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/book/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Book not found'. Got '%s'", m["error"])
	}
}

func TestCreateBook(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test book", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test book" {
		t.Errorf("Expected book name to be 'test book'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected book price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected book ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetBook(t *testing.T) {
	clearTable()
	addBook(1)

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addBook(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO book(name, price) VALUES($1, $2)", "Book "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateBook(t *testing.T) {

	clearTable()
	addBook(1)

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	var jsonStr = []byte(`{"name":"test book - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/book/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalBook["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalBook["id"], m["id"])
	}

	if m["name"] == originalBook["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalBook["name"], m["name"], m["name"])
	}

	if m["price"] == originalBook["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalBook["price"], m["price"], m["price"])
	}
}

func TestDeleteBook(t *testing.T) {
	clearTable()
	addBook(1)

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/book/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/book/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
