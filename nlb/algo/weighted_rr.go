package algo

import (
	"errors"
	"math/rand"
	"nlb/k8s"
)

type WeightedRoundrobin struct {
	Index int
}

func (rr *WeightedRoundrobin) GetIP(ips *[]*k8s.PodDetails, clientAddress string) (string, error) {

	ip_lst := *ips
	var cpu []float64
	var serverIP []string
	var total float64

	for _, ip := range ip_lst {
		serverIP = append(serverIP, ip.IP)
		if len(cpu) == 0 {
			cpu = append(cpu, ip.Memory)
		} else {
			cpu = append(cpu, ip.Memory+cpu[len(cpu)-1])
		}
		total += ip.Memory
	}

	randNumber := (rand.Float64()) * total
	for index, cdf := range cpu {
		if randNumber < cdf {
			return serverIP[index], nil
		}
	}

	return "", errors.New("failed to get ip")
}
