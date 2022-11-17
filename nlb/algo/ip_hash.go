package algo

import (
	//"fmt"
	"log"
	"strings"
	"strconv"
	"nlb/k8s"
)

type Ip_Hash struct{
	Ip string
	Port string
}

func hash(ip string, port string, n int, ips []*k8s.PodDetails)(string,error){

	//ip_lst := *ips

	//parse string by .
	num_list := strings.FieldsFunc(ip, func(r rune) bool {
		if r == '.' {
			return true
		}
		return false
	})

	total_num, err := strconv.Atoi(port)
	if err != nil{
		return "", err
	}
	//convert string to ints, add ints
	for _,num := range num_list {
		new_num,err := strconv.Atoi(num)
		if err != nil{
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

func (ih *Ip_Hash) GetIP(ips *[]*k8s.PodDetails)(string, error){
	ip_lst := *ips

	log.Print("hashing ip: ", ih.Ip)
	log.Print("hashing port: ", ih.Port)

	ip,err := hash(ih.Ip, ih.Port, len(ip_lst), ip_lst)

	return ip, err
}
