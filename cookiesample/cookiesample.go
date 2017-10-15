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
 const (
 	cookieName = "test"
 	header = "Set-Cookie"
 )

var (
	key = []byte("the-key-has-to-be-32-bytes-long!")


)


// This is an example of how to encrypt a struct and attach it to a cookie


type UserRights struct {
	username string   	// This is the user's name
	// define the user's roles
	admin    bool		// These are the roles - this user has admin rights
	edit     bool       // The user has edit rights
	user     bool		// This is an ordinary user
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
		if err != nil {
			// okay handle this - redirect the user to an error page
			fmt.Fprintf(w, "Error creating the cookie, %s ",err.Error())
			return
		}

		// Create the cookie
		maxAge := 8 * 60 * 60 // expires after 8 hours
		cookie := http.Cookie{
			Name:  cookieName,
			Value: base64.StdEncoding.EncodeToString(cookieValue),
			MaxAge: maxAge,
			Secure: true,
		}

		// Store the cookie in the response to send it back to the browser
		http.SetCookie(w, &cookie)
		fmt.Fprintf(w, "Login Okay")
		return
	}

	fmt.Fprintf(w, "Invalid User")
}

// Validation is beyond the scope to this sample
// Always validate any user
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
		user: true,
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


// This is a handler func - decide what to server the user based on roles
func SomeHandlerFunc(w http.ResponseWriter, r *http.Request) {

	userRights, err := validateRequest(&w,r)
	if err == http.ErrNoCookie {
		// Display the login page - the cookie doesn't exist
		fmt.Fprintf(w, "Please Login")
	} else if err  != nil {
		fmt.Fprintf(w, "Error %s",err.Error())
		// display an error - this is an invalid cookie
	}

	// OKAY at this point diplay the page + data the user's after
	// Check the rights and display whatever data is necessary eg:
	if userRights.admin {
		// This is an admin user display stuff
		fmt.Fprintf(w, "Hi admin user")
	} else if userRights.edit {
		// This user has edit rights
		fmt.Fprintf(w, "Hi edit user")
	} else if userRights.user {
		// This user is a normal user
		fmt.Fprintf(w, "Hi ordinary user")
	} else {
		// This is an invalid user struct - display an error
	}
}

// This is any old handler func - server all users the same thing
func SomeOtherHandlerFunc(w http.ResponseWriter, r *http.Request) {

	_, err := validateRequest(&w,r)
	if err == http.ErrNoCookie {
		// Display the login page - the cookie doesn't exist
		fmt.Fprintf(w, "Please login")
		return
	} else if err  != nil {
		// display an error - this is an invalid cookie
		fmt.Fprintf(w, "Error... %s",err.Error())
		return
	}

	// This is an admin user display stuff
	fmt.Fprintf(w, "Hi user")
}


// Validate the cookie. Is the user logged in? What roles does he/she have?
// This is a key function - apply this to all handlers
func validateRequest(w * http.ResponseWriter, r *http.Request) (*UserRights, error) {

	cookie, err := r.Cookie(cookieName)
	if err  != nil {
		return nil, err
	}

	encoded, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		// Cannot decode the cookie - ensure that the handler displays the error page
		return nil, err
	}

	rights, err := decrypt(encoded,key)
	return rights, nil
}


// Decrypt the bytes, converting them back to to a struct
func decrypt(ciphertext []byte, key []byte) (*UserRights, error) {
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
	var decoded []byte
	decoded, err = gcm.Open(nil, nonce, ciphertext, nil)

	var userRights UserRights
	err = json.Unmarshal(decoded,&userRights)

	if err != nil {
		// This should lead to the display of an error page
		fmt.Printf("Error: %s", err)
		return nil, err;
	}

	return &userRights,nil
}


