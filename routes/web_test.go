package routes

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"miniblog/domain/mocks"
	"miniblog/domain/service"
	"miniblog/domain/validator"
	c "miniblog/http/controller"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var router *mux.Router

// move dir to root to load templates successfully
func init() {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	//DB initialize
	println("-------setup-------")
	router = InitialRouters()
	// User
	userService := &service.UserService{}
	userHandler := &c.UserHandler{Service: userService}
	SetUserRouters(router, userHandler)

	//Article
	articleRepository := new(mocks.ArticleRepository)
	articleService := &service.ArticleService{Repo: articleRepository, V: validator.NewValidator()}
	articleHandler := &c.ArticleHandler{Service: articleService, UserService:userService}
	SetArticleRouters(router, articleHandler)
}

// mock request and recorder by httptest
func TestHealthCheckHandler(t *testing.T) {
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

// TODO: research other good method to check handler
func TestHandlerStatus(t *testing.T) {
	var tests = []struct {
		in  string
		isLogin bool
		out int
	}{
		{"/login", false, http.StatusOK},
		{"/logout", false, http.StatusSeeOther},
		{"/", false, http.StatusOK},
		{"/article/1", false, http.StatusOK},
		{"/articles/new", true, http.StatusOK},
		{"/articles/new", false, http.StatusFound},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.in, nil)
			w := httptest.NewRecorder()
			if tt.isLogin == true {
				http.SetCookie(w, &http.Cookie{Name: "user", Value: "demo"})
				r.Header.Set("Cookie", w.Header()["Set-Cookie"][0])
			}
			router.ServeHTTP(w, r)
			if w.Code != tt.out {
				t.Errorf("expected status OK, got %v", w.Code)
			}
		})
	}
}

func TestHandlerStatusPostForm(t *testing.T) {
	var tests = []struct {
		in  string
		body url.Values
		isLogin bool
		out int
	}{
		{"/login", url.Values{"name":{"demo"}, "password":{"demo"}}, false, http.StatusSeeOther},
		{"/login", url.Values{"name":{"abc"}, "password":{"ddd"}}, false, http.StatusOK},
		{"/login", url.Values{}, false, http.StatusOK},
		{"/articles", url.Values{"title":{"test title"}, "description":{"test description"}}, true, http.StatusSeeOther},
		{"/articles", url.Values{}, true, http.StatusUnprocessableEntity},
		{"/articles", url.Values{}, false, http.StatusFound},

	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", tt.in, strings.NewReader(tt.body.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if tt.isLogin == true {
				http.SetCookie(w, &http.Cookie{Name: "user", Value: "demo"})
				r.Header.Set("Cookie", w.Header()["Set-Cookie"][0])
			}

			router.ServeHTTP(w, r)
			if w.Code != tt.out {
				t.Errorf("expected status OK, got %v", w.Code)
			}
		})
	}
}

// launch server and test routing
func TestRoutingWithLaunchedServer(t *testing.T) {
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

