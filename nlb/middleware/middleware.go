package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/pbkdf2"
)

var (
	salt       = "saltysalt"
	iv         = "                "
	length     = 16
	password   = ""
	iterations = 1003
)

// The server can set a cookie
func SetCookie(w http.ResponseWriter, req *http.Request, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "nlb-cookie_abcde",
		Value: value,
		Path:  "/",
	})
	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(w, "in chrome go to: dev tools / application / cookies")
}

func ReadCookie(w http.ResponseWriter, req *http.Request) string {
	cookie, err := req.Cookie("nlb-cookie_abcde")
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

func EncryptMessage(key string, message string) string {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(message))
	fmt.Println("Message length", len([]byte(message)))
	c.Encrypt(msgByte, []byte(message))
	return hex.EncodeToString(msgByte)
}

// func DecryptMessage(key string, message string) string {
// 	txt, _ := hex.DecodeString(message)
// 	c, err := aes.NewCipher([]byte(key))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	msgByte := make([]byte, len(txt))
// 	c.Decrypt(msgByte, []byte(txt))

//		msg := string(msgByte[:])
//		return msg
//	}

// https://gist.github.com/dacort/bd6a5116224c594b14db
func DecryptValue(encryptedValue []byte) string {
	key := pbkdf2.Key([]byte(password), []byte(salt), iterations, length, sha1.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	decrypted := make([]byte, len(encryptedValue))
	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	cbc.CryptBlocks(decrypted, encryptedValue)

	plainText, err := aesStripPadding(decrypted)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return ""
	}

	return string(plainText)
}

// In the padding scheme the last <padding length> bytes
// have a value equal to the padding length, always in (1,16]
func aesStripPadding(data []byte) ([]byte, error) {
	if len(data)%length != 0 {
		return nil, fmt.Errorf("decrypted data block length is not a multiple of %d", length)
	}
	paddingLen := int(data[len(data)-1])
	if paddingLen > 16 {
		return nil, fmt.Errorf("invalid last block padding length: %d", paddingLen)
	}
	return data[:len(data)-paddingLen], nil
}
