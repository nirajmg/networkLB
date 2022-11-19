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
	"time"
)

var Ips *[]*k8s.PodDetails
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
		cookies := r.Cookies()
		serverIp := ""
		// ipEncrypt := ""
		isCookieExist := middleware.CookieExists(cookies, "nlb-cookie_abcde")

		if isCookieExist {
			fmt.Println("HERE!: Cookie exists!")
			serverIp = middleware.ReadCookie(w, r)
			// byteStr := []byte(serverIp)
			// fmt.Println("Decrypt:", middleware.DecryptValue(byteStr))
			// decryptedMessage := string(middleware.DecryptMessage("nlb-cookie_abcde", encryptedIp))
			// strArr := strings.Split(decryptedMessage, "_")
			// serverIp = strArr[0]
			//TODO: Strip the cookie information (LATER)
		} else {
			fmt.Println("HERE!: Client does not have a cookie, generating...")
			//Get a random ip and set serverIp to the random server ip
			serverIp, _ = algoIP.GetIP(Ips)
			//Encrypt the server ip and set that as the value of the cookie
			// ipEncrypt = middleware.EncryptMessage("nlb-cookie_abcde", serverIp+"_abcdef")
		}
		fmt.Println("Current IP: ", serverIp)
		fmt.Println("Ips: ", Ips)

		algoIP = &algo.Ip_Hash{Ip: "192.168.0.1", Port: "8000"}
		serverIp, _ = algoIP.GetIP(Ips)

		proxy, err := NewProxy("http://" + serverIp + ":80") //change this line
		if err != nil {
			panic(err)
		}

		//Configuration here to server, if we get a statuscode of 200 then set the cookie for the client
		if isCookieExist {
			fmt.Println("Setting cookie... ", serverIp)
			proxy.ModifyResponse = func(res *http.Response) error {
				if res.StatusCode == 200 {
					//Set cookie for the client
					middleware.SetCookie(w, r, serverIp)
				}
				return nil
			}
		}
		proxy.ServeHTTP(w, r)
	}
}

func UpdateIP() {
	for {
		ips, err := k8s.ListPod()
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%v", ips)
		Ips = &ips
		time.Sleep(2 * time.Second)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	print("In NLB\n")
	if err := k8s.NewClient(); err != nil {
		panic(err)
	}

	algoIP = &algo.WeightedRoundrobin{Index: 0}

	go func() {
		UpdateIP()
	}()

	time.Sleep(2 * time.Second)
	algoIP.GetIP(Ips)

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler())
	http.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(":80", nil))
}
