package algo

import (
	//"fmt"
	"fmt"
	"log"
	"nlb/k8s"
	"regexp"
	"strconv"
	"strings"
)

type Ip_Hash struct {
	Address string
}

func hash(ip string, port string, n int, ips []*k8s.PodDetails) (string, error) {

	//ip_lst := *ips

	//parse string by .
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
	log.Print("sum: ", total_num)
	// % by n
	index := total_num % n
	log.Print("index: ", index)
	//lookup corresponding ip
	server_ip := ips[index].IP
	log.Print("server ip: ", server_ip)

	return server_ip, err
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func (ih *Ip_Hash) GetIP(ips *[]*k8s.PodDetails) (string, error) {
	ip_lst := *ips

	s := strings.Split(ih.Address, ":")
	fmt.Println(s)
	log.Print("hashing ip: ", s[0])
	log.Print("hashing port: ", s[1])

	if !validIP4(s[0]) {
		s[0] = "127.0.0.1"
	}
	ip, err := hash(s[0], s[1], len(ip_lst), ip_lst)

	return ip, err
}
