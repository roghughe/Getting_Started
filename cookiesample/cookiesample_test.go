package cookiesample

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"
	"bufio"
	"strings"
)

func TestLogin(t *testing.T) {
	fmt.Println("TestLogin - - Running")

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	Login(w, req)

	headerMap := w.HeaderMap
	cookie := headerMap.Get("Set-Cookie")
	fmt.Printf("Cookie is: %+v\n",cookie)
	result := string(w.Body.Bytes())
	if result != "Login Okay" {
		t.Fatalf("Invalid body: %s\n",result)
	}

}



func TestSomeOtherHandlerFunc(t *testing.T) {

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	Login(w, req) // The user must be logged in

	// Copy the Cookie over to a new Request
	// The browser will do this bit before it makes the next requrest
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	rawCookies := w.Header().Get("Set-Cookie")
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookies)

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	var cookie *http.Cookie
	if err == nil {
		cookie, err = req.Cookie(cookieName)
		req.AddCookie(cookie)

		// Now ensure that we get the page we want
		SomeOtherHandlerFunc(w,req)

		body := w.Body
		result := string(body.Bytes())
		if result != "Login OkayHi user" {
			t.Fatalf("Invalid body: %s\n",result)
		}
		return
	}

	t.Error("Can't copy cookie to the new request")
}


func TestSomeOtherHandlerFunc_user_not_logged_in(t *testing.T) {

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// Login not called
	SomeOtherHandlerFunc(w,req)

	body := w.Body

	result := string(body.Bytes())

	if result != "Please login" {
		t.Fatalf("Invalid body: %s\n",result)
	}
}
