// Source: https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"nlb/algo"
	"nlb/k8s"
	"nlb/middleware"
	"strings"
	"time"
)

var Ips *[]string
var algoIP algo.Algorithm

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	fmt.Println("URL: ", url)
	p := httputil.NewSingleHostReverseProxy(url)
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
		// generate a encryption, map the encryption to ip

		cookies := r.Cookies()
		fmt.Println("Cookies: ", cookies)
		serverIp := ""
		ipEncrypt := ""
		isCookieExist := middleware.CookieExists(cookies, "nlb-cookie_abcde")

		if isCookieExist {
			encryptedIp := middleware.Read(w, r)
			decryptedMessage := string(middleware.DecryptMessage("nlb-cookie_abcde", encryptedIp))
			strArr := strings.Split(decryptedMessage, "_")
			serverIp = strArr[0]
			//TODO: Strip the cookie information (LATER)
		} else {
			fmt.Println("Client does not have a cookie, generating...")
			//Get a random ip and set serverIp to the random server ip
			serverIp, _ = algoIP.GetIP(Ips)
			//TODO: Maybe use a hash to generate the message for encryption
			//Encrypt the server ip and set that as the value of the cookie
			ipEncrypt = middleware.EncryptMessage("nlb-cookie_abcde", serverIp+"_abcdef")
		}

		proxy, err := NewProxy("http://" + serverIp + ":80") //change this line
		if err != nil {
			panic(err)
		}
		//Configuration here to server, if we get a statuscode of 200 then set the cookie for the client
		proxy.ModifyResponse = func(res *http.Response) error {
			if res.StatusCode == 200 {
				//Set cookie for the client
				fmt.Println("Encrypted ip:", ipEncrypt)
				middleware.Set(w, r, ipEncrypt)
			}
			return nil
		}
		proxy.ServeHTTP(w, r)
	}
}

func UpdateIP() {

	ips, err := k8s.ListPod()
	if err != nil {
		panic(err)
	}

	fmt.Println(ips)
	Ips = &ips
	time.Sleep(2 * time.Second)

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
	algoIP = &algo.Roundrobin{Index: 0}

	go func() {
		UpdateIP()
	}()
	time.Sleep(2 * time.Second)
	ip, _ := algoIP.GetIP(Ips)
	print(ip)

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler())
	http.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(":80", nil))
}
