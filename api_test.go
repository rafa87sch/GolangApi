package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var limit = 65
var ErrTooLarge = errors.New("String is too large!")

//testing metho
func TestCreateEntryEncrypt(t *testing.T) {

	var jsonStr = []byte(`{"Password":"123"}`)

	req, err := http.NewRequest("POST", "api/encrypt", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(encrypt)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if err := Process(rr.Body.String()); err != nil {
		t.Errorf("The encrypt is not working since the length is not matching %v", err)
	}
}
func TestCreateEntryDecrypt(t *testing.T) {

	var jsonStr = []byte(`{"Password":"4521ef5e25909739a43938d60064ebf6159e3ab90183d9c5d859435f599cf4"}`)

	req, err := http.NewRequest("POST", "/api/decrypt", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(decrypt)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `"123"`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//taking the len from the encry
func Process(s string) error {
	if len(s) != limit {
		return ErrTooLarge
	}

	// All OK
	return nil
}
