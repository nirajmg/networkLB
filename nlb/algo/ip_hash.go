package algo

import (
	//"fmt"
	"strings"
	"strconv"
	"container/list"
)

type Ip_Hash struct{
	Ip string
	Port string
}

func hash(ip string, port string, n int, ips *[]string)(string,error){

	//parse string by .
	num_list := strings.FieldsFunc(ip, func(r rune) bool {
		if r == '.' {
			return true
		}
		return false
	})
	total_num int := strconv.Atoi(port)
	//convert string to ints, add ints
	for _,num := range ips {
		new_num,err := strconv.Atoi(num)
		if err != nil{
			return nil, err
		}
		total_num += new_num
	}
	// % by n
	index := total_num % n
	//lookup corresponding ip
	server_ip := ips[index]

	return server_ip, err
}

func (ih *Iphash) Get_IP(ips *[]string)(string, error){
	ip_lst := *ips
	if ih.Index >= len(ip_lst) {
		ih.Index = 0
	}
	ip,err := hash(ih.IP, len(ip_list), *ip_lst)

	return ip, err
}
