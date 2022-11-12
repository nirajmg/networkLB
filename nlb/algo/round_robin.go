package algo

import "nlb/k8s"

type Roundrobin struct {
	Index int
}

func (rr *Roundrobin) GetIP(ips *[]*k8s.PodDetails) (string, error) {
	ip_lst := *ips
	if rr.Index >= len(ip_lst) {
		rr.Index = 0
	}

	ip := ip_lst[rr.Index]
	rr.Index += 1

	return ip.IP, nil

}
