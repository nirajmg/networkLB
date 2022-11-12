package algo

import (
	"net/http"
	"nlb/k8s"
	"time"
)

type LeastResTime struct {
	ResTimes []int
}

func (lrt *LeastResTime) GetIP(ips *[]*k8s.PodDetails) (string, error) {

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	ip_lst := *ips

	for _, ip := range ip_lst {

		start := time.Now()
		client.Get("https://" + ip.IP + ":80/")
		elapsed := time.Now().Sub(start)
		lrt.ResTimes = append(lrt.ResTimes, int(elapsed))
	}

	var m, minIdx = 0, 0
	for i, e := range lrt.ResTimes {
		if i == 0 || e < m {
			m = e
			minIdx = i
		}
	}
	return ip_lst[minIdx].IP, nil
}
