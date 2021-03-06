package apiserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	result, err := Hash("hello")
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	expected := "5d41402abc4b2a76b9719d911017c592"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(ping).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected %d, got %d", 200, status)
	}

	respBody := rr.Body.String()
	if respBody != "Pong" {
		t.Errorf("Expected Pong, got %s", respBody)
	}

}

func TestShortern(t *testing.T) {
	api := NewServer()
	reader := strings.NewReader(`{"URL":"https://github.com/go-redis/redis"}`)
	req, err := http.NewRequest("POST", "/shortern", reader)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(api.shotern).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected %d, got %d", 200, status)
	}

	expectedResponse := string(`{"message":"Success","data":"f7c126d0514c781a6947d90b37e384c2"}`)
	respBody := rr.Body.String()
	assert.JSONEq(t, expectedResponse, respBody)

}
