package cookiesample

import (
	"testing"
	"fmt"
	"net/http/httptest"
)

func TestLogin(t *testing.T) {
	fmt.Println("TestLogin - - Running")

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	Login(w, req)

	headerMap := w.HeaderMap
	cookie := headerMap.Get("Set-Cookie")

	fmt.Printf("Cookie is: %+v\n",cookie)
}

