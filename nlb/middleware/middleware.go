package middleware

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"net/http"
)

// The server can set a cookie
func Set(w http.ResponseWriter, req *http.Request, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "my-cookie",
		Value: value,
		Path:  "/",
	})
	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(w, "in chrome go to: dev tools / application / cookies")
}

func Read(w http.ResponseWriter, req *http.Request) string {
	cookie, err := req.Cookie("my-cookie")
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return ""
	}
	return cookie.String()
}

func Expire(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		return
	}
	c.MaxAge = -1 // delete cookie
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func CookieExists(cookies []*http.Cookie, cookieName string) bool {
	if len(cookies) == 0 {
		return false
	} else {
		for _, c := range cookies {
			if c.Name == cookieName {
				return true
			}
		}
		return false
	}
}

func EncryptMessage(key string, message string) string {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(message))
	c.Encrypt(msgByte, []byte(message))
	return hex.EncodeToString(msgByte)
}

func DecryptMessage(key string, message string) string {
	txt, _ := hex.DecodeString(message)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(txt))
	c.Decrypt(msgByte, []byte(txt))

	msg := string(msgByte[:])
	return msg
}
