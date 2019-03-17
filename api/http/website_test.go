package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
)

func TestWebsiteHandler(t *testing.T) {
	h := NewTestHandler()
	req, err := http.NewRequest("GET", "/api/websites", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.GET("/api/websites", h.Websites)

	router.ServeHTTP(rr, req)
	spew.Dump(rr)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
