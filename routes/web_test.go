package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock request and recorder by httptest
func TestHealthCheckHandler(t *testing.T) {
	router := InitialRouters()

	// we don't check r & w because we want to check handler functions
	r := httptest.NewRequest("GET", "/health_check", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	// check handler function
	// check http status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}

	// check response
	respString := w.Body.String()
	expected := "I'm fine!"

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

// launch server and test routing
func TestRoutingWithLaunchedServer(t *testing.T) {
	router := InitialRouters()

	// Create a new mock server
	mockServer := httptest.NewServer(router)
	defer mockServer.Close()

	res, err := http.Get(mockServer.URL + "/health_check")
	if err != nil {
		t.Errorf("could not send Get request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}

	// read the body into a bunch of bytes (b)
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("could not read body: %v", err)
	}
	// convert the bytes to a string
	respString := string(b)
	expected := "I'm fine!"

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

