package cookiesample

import (
	"crypto/aes"

	"encoding/json"
	"encoding/base64"
	"fmt"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"net/http"
	"errors"
)

var key = []byte("the-key-has-to-be-32-bytes-long!")



// This is an example of how to encrypt a struct and attach it to a cookie


type UserRights struct {
	username string   	// This is the user's name
	// define the user's roles
	admin    bool		// These are the roles - this user has admin rights
	edit     bool       // The user has edit rights
	tracking bool
}



// This is the sample handler function
func Login(w http.ResponseWriter, r *http.Request) {

	// Get the username and password from the request
	form  := r.Form
	username := form.Get("username")
	password := form.Get("password")

	// validate the password
	if validateUser(username,password) {

		// Get the user's rights
		rights := getRights(username)
		cookieValue, err := encrypt(rights)
		if err == nil {
			// okay handle this - redirect the user to an error page
		}

		// Create the cookie
		maxAge := 8 * 60 * 60 // expires after 8 hours
		cookie := http.Cookie{
			Name:  "test",
			Value: base64.StdEncoding.EncodeToString(cookieValue),
			MaxAge: maxAge,
			Secure: true,
		}

		// Store the cookie in the response to send it back to the browser
		v := cookie.String()
		w.Header().Add("Set-Cookie", v)
	}
}

// Validation is beyond the scope to this sample
func validateUser(username, password string) bool {
	return true
}

// Figure out what this user has access to
func getRights(username string) *UserRights {

	// This may go to a database here
	// In this case they're fudged
	return &UserRights{
		username: username,
		admin: true,
		edit:  true,
		tracking: true,
	}
}

// This is boilerplate code
func encrypt(rights *UserRights) ([]byte, error) {

	b, err := json.Marshal(rights)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err;
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, b, nil), nil
}


// This is any old handler func - this does anything, but not before decrypting the cookie
func someHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if ok, err := validateRequest(&w,r); err != nil {

		if err != nil {
			// display an error - this is an invalid cookie
		}

		if !ok {
			// At this point display 'invalid user page'
		}
	}


	// OKAY at this point diplay the page + data the user's after


}

func validateRequest(w * http.ResponseWriter, r *http.Request) (bool, error) {


	cookie, err := r.Cookie("test")
	if err != nil {
		fmt.Printf("Cookie error: %+v\n", err.Error())
		// Okay, so return an error page
	}
	return true, nil
}




func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}


