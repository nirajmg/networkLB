// Source: https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"nlb/k8s"
)

// The server can set a cookie
func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "my-cookie",
		Value: "some value",
		Path:  "/",
	})
	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(w, "in chrome go to: dev tools / application / cookies")
}

func read(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("my-cookie")
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "YOUR COOKIE:", c)
}

func expire(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		return
	}
	c.MaxAge = -1 // delete cookie
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	print("Here we are!\n")
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		fmt.Printf("%d\n", len(cookies))
		if len(cookies) == 0 {
			set(w, r)
		} else {
			for _, c := range cookies {
				if c.Name == "my-cookie" {
					read(w, r)
				} else {
					set(w, r)
				}
			}
		}
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("http://google.com")
	if err != nil {
		panic(err)
	}

	if err := k8s.NewClient(); err != nil {
		panic(err)
	}

	ip, err := k8s.GetPodDetails("postgresql-0")
	if err != nil {
		panic(err)
	}
	print(ip)

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
