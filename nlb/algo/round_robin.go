package algo

type Roundrobin struct {
	index int
}

func (rr *Roundrobin) GetIP(ip_lst []string) (string, error) {
	if rr.index >= len(ip_lst) {
		rr.index = 0
	}

	ip := ip_lst[rr.index]
	rr.index += 1

	return ip, nil

}

func New() *Roundrobin {
	return &Roundrobin{
		index: 0,
	}
}
