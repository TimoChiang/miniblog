package service

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const userName = "test"

func TestUserService_GetUser(t *testing.T) {
	// request without cookie
	r1 := httptest.NewRequest("GET", "/", nil)

	// request with cookie
	w := httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: userName})
	r2 := &http.Request{Header: http.Header{"Cookie": w.Header()["Set-Cookie"]}}

	var tests = []struct {
		in  string
		r  *http.Request
		out *User
	}{
		{"request without cookie", r1, nil},
		{"request with cookie",r2, &User{userName}},
	}
	service := &UserService{}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			user := service.GetUser(tt.r)
			if !reflect.DeepEqual(user, tt.out) {
				t.Errorf("got %q, want %q", user, tt.out)
			}
		})
	}
}

func TestUserService_SetCookie(t *testing.T) {
	w := httptest.NewRecorder()
	service := &UserService{}
	service.SetCookie(w)
	cookie := w.Result().Cookies()[0]
	expected := &http.Cookie{
		Name:     cookieName,
		Value:    "demo",
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
		MaxAge: 86400}

	// TODO: compare cookie directly
	//if !reflect.DeepEqual(cookie, expected) {
	//	t.Errorf("got %q, want %q", cookie, expected)
	//}
	if cookie.Name != expected.Name {
		t.Errorf("cookie Name got %q, want %q", cookie.Name, expected.Name)
	}

	if cookie.Value != expected.Value {
		t.Errorf("cookie Value got %q, want %q", cookie.Value, expected.Value)
	}

	if cookie.Path != expected.Path {
		t.Errorf("cookie Path got %q, want %q", cookie.Path, expected.Path)
	}

	if cookie.HttpOnly != expected.HttpOnly {
		t.Errorf("cookie HttpOnly got %t, want %t", cookie.HttpOnly, expected.HttpOnly)
	}

	if cookie.MaxAge != expected.MaxAge {
		t.Errorf("cookie MaxAge got %q, want %q", cookie.MaxAge, expected.MaxAge)
	}
}

func TestUserService_DeleteCookie(t *testing.T) {
	// initialize cookie
	w := httptest.NewRecorder()
	//http.SetCookie(w, &http.Cookie{Name: cookieName, Value: userName})

	service := &UserService{}
	service.DeleteCookie(w)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %q, want no cookie", cookies)
	}

	cookie := cookies[0]
	expected := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
		MaxAge: 0}

	// TODO: compare cookie directly
	//if !reflect.DeepEqual(cookie, expected) {
	//	t.Errorf("got %q, want %q", cookie, expected)
	//}
	if cookie.Name != expected.Name {
		t.Errorf("cookie Name got %q, want %q", cookie.Name, expected.Name)
	}

	if cookie.Value != expected.Value {
		t.Errorf("cookie Value got %q, want %q", cookie.Value, expected.Value)
	}

	if cookie.Path != expected.Path {
		t.Errorf("cookie Path got %q, want %q", cookie.Path, expected.Path)
	}

	if cookie.HttpOnly != expected.HttpOnly {
		t.Errorf("cookie HttpOnly got %t, want %t", cookie.HttpOnly, expected.HttpOnly)
	}

	if cookie.MaxAge != expected.MaxAge {
		t.Errorf("cookie MaxAge got %q, want %q", cookie.MaxAge, expected.MaxAge)
	}

	now := time.Now()
	if cookie.Expires.After(now) {
		t.Errorf("cookie Expires is %q, want earlier then %q", cookie.Expires, now)
	}
}
