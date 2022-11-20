package middleware

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zellyn/kooky/browser/chrome"
)

var CookieKey = "nlb-cookie"

func SetCookie(w http.ResponseWriter, req *http.Request, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  CookieKey,
		Value: value,
		Path:  "/",
	})

	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(w, "in chrome go to: dev tools / application / cookies")
}

func ReadCookie(w http.ResponseWriter, req *http.Request) string {
	cookie, err := req.Cookie(CookieKey)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
	}
	return cookie.Value
}

func ExpireCookie(w http.ResponseWriter, req *http.Request) {
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

func EncryptMessage(key string, message string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	msgByte := make([]byte, len(message))
	fmt.Println("Message length", len([]byte(message)))
	c.Encrypt(msgByte, []byte(message))
	return hex.EncodeToString(msgByte), nil
}

func DecryptMessage(key string, message string) (string, error) {
	txt, _ := hex.DecodeString(message)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	msgByte := make([]byte, len(txt))
	c.Decrypt(msgByte, []byte(txt))

	return string(msgByte[:]), nil
}

func GetCookie() {
	dir, _ := os.UserConfigDir() // "/<USER>/Library/Application Support/"
	cookiesFile := dir + "/Google/Chrome/Default/Cookies"
	cookies, err := chrome.ReadCookies(cookiesFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, cookie := range cookies {
		fmt.Println(cookie)
	}
}
