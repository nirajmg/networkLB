package algo

import (
	"nlb/k8s"
	"regexp"
	"strconv"
	"strings"
)

type Ip_Hash struct{}

func hash(ip string, port string, ips []*k8s.PodDetails) (string, error) {

	num_list := strings.FieldsFunc(ip, func(r rune) bool {
		if r == '.' {
			return true
		}
		return false
	})

	total_num, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}
	//convert string to ints, add ints
	for _, num := range num_list {
		new_num, err := strconv.Atoi(num)
		if err != nil {
			return "", err
		}
		total_num += new_num
	}
	index := total_num % len(ips)

	return ips[index].IP, err
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	ipRegex := `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	re, _ := regexp.Compile(ipRegex)
	return re.MatchString(ipAddress)
}

func (ih *Ip_Hash) GetIP(ips *[]*k8s.PodDetails, clientAddress string) (string, error) {
	s := strings.Split(clientAddress, ":")
	if !validIP4(s[0]) {
		s[0] = "127.0.0.1"
	}

	return hash(s[0], s[1], *ips)
}
