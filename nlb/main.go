// Source: https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"nlb/k8s"
	"time"
)

var ips []string

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

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
	p := httputil.NewSingleHostReverseProxy(url)
	// p.Director = func(w *http.Response) {
	// 	w.Header.Set("cookie", "shit")
	// }
	//How to modify responses
	p.ModifyResponse = func(res *http.Response) error {
		if res.StatusCode == 200 {
			res.Header.Set("cookie", "cook")
		}
		return nil
	}
	return p, nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler() func(http.ResponseWriter, *http.Request) {
	//parse the request
	fmt.Println("In Proxy Request Handler")
	return func(w http.ResponseWriter, r *http.Request) {
		// angel has a cookie
		// read the cookie , remove from the request , take the ip and send it to that ip
		// just need r
		// also make sure that return has cookie

		// oh no she doesnt
		// generate cookie for now put a random server ip cookie = ip
		// this should happend after server finished the response
		// generate a hash, map the hash to ip
		fmt.Println("Cookies")
		cookies := r.Cookies()
		fmt.Printf("%d\n", cookies)

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

		proxy, err := NewProxy("http://localhost:30691") //change this line
		if err != nil {
			panic(err)
		}
		//stripping the cookie information

		//Set cookie for the client
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	print("In NLB\n")
	if err := k8s.NewClient(); err != nil {
		panic(err)
	}

	_, err := k8s.GetPodDetails("postgresql-0")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			ips, err := k8s.ListPod()
			if err != nil {
				panic(err)
			}
			fmt.Println(ips)
		}
	}()

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler())
	http.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(":80", nil))
}
