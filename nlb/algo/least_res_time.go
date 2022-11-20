package algo

import (
	"fmt"
	"net/http"
	"nlb/k8s"
	"time"
)

type LeastResTime struct {
	ResTimes []int
}

func (lrt *LeastResTime) GetIP(ips *[]*k8s.PodDetails, clientAddress string) (string, error) {

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	ip_lst := *ips

	for _, ip := range ip_lst {
		start := time.Now()
		client.Get(fmt.Sprintf("http://%s:80/", ip.IP))
		elapsed := time.Now().Sub(start)
		lrt.ResTimes = append(lrt.ResTimes, int(elapsed))
	}

	var minTime, minIdx = lrt.ResTimes[0], 0
	for index, resTime := range lrt.ResTimes {
		if index == 0 || resTime < minTime {
			minTime = resTime
			minIdx = index
		}
	}
	return ip_lst[minIdx].IP, nil
}
