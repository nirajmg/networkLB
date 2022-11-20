// Source: https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/
package main

import (
	"fmt"

	"net/http"
	"net/http/httputil"
	"net/url"
	"nlb/algo"
	"nlb/k8s"
	"nlb/middleware"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var Ips *[]*k8s.PodDetails
var algoIP algo.Algorithm

func serverStats() {
	for {
		ips, err := k8s.ListPod()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to Update server details")
			continue
		}
		Ips = &ips
		time.Sleep(2 * time.Second)
	}
}

func setLBAlgorithm() {
	algoId, exists := os.LookupEnv("LB_ALGO")
	if !exists {
		log.Error("LB_ALGO env missing defaulting to round robin")
		algoId = "rr"
	}
	switch algoId {
	case "rr":
		log.Info("Load balancing using round robin")
		algoIP = &algo.Roundrobin{}
	case "iph":
		log.Info("Load balancing using ip hash")
		algoIP = &algo.Ip_Hash{}
	case "wrr":
		log.Info("Load balancing using weighted round robin")
		algoIP = &algo.WeightedRoundrobin{}
	case "lrt":
		log.Info("Load balancing using least response time")
		algoIP = &algo.LeastResTime{}
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler() func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		serverIp := ""
		var err error
		var IPTime string

		cookieExists := middleware.CookieExists(cookies, middleware.CookieKey)
		if cookieExists {
			serverIp = middleware.ReadCookie(w, r)
		} else {
			start := time.Now()
			serverIp, err = algoIP.GetIP(Ips, r.RemoteAddr)
			IPTime = time.Now().Sub(start).String()
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Fall back to Round Robin")
				fallBackAlgo := &algo.Roundrobin{}
				serverIp, err = fallBackAlgo.GetIP(Ips, r.RemoteAddr)
			}
		}

		log.WithFields(log.Fields{
			"client": r.RemoteAddr,
			"server": serverIp,
			"cookie": cookieExists}).Info("Request Details")

		proxy, err := NewProxy(fmt.Sprintf("http://%s:80/", serverIp)) //change this line
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to connect to backend")
			return
		}

		if !cookieExists {
			proxy.ModifyResponse = func(res *http.Response) error {
				w.Header().Add("X-Metrics-IP", serverIp)
				w.Header().Add("X-Metrics-Time", IPTime)
				if res.StatusCode == 200 {
					middleware.SetCookie(w, r, serverIp)
				}
				return nil
			}
		}

		proxy.ServeHTTP(w, r)
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	log.Info("Initializing kubernetes client")
	if err := k8s.NewClient(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatalf("Failed to create a k8s client")
	}

	log.Info("Starting routine for server updates")
	go func() {
		serverStats()
	}()

	log.Info("Setting up the LB algorithms")
	setLBAlgorithm()

	log.Info("Starting the server")

	http.HandleFunc("/", ProxyRequestHandler())
	http.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(":80", nil))
}
