package main

import (
	"bufio"
	"fmt"
	"github.com/AlekSi/zabbix"
	"github.com/howeyc/gopass"
	"os"
	"strings"
)

var (
	serverSelected string
)

func main() {
	// Login to the server
	api := Login()

	DeleteHosts(api)
}

func Ask4Confirm() bool {
	var s string

	fmt.Printf("(y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func GetCreds() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter zabbix username: ")
	u, _ := reader.ReadString('\n')
	u = strings.TrimSpace(u)
	fmt.Printf("Password: ")
	pByte, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p := string(pByte)
	return u, p
}

func Login() zabbix.API {

	serverA := "zabbix-dc1.yourcompany.com"
	serverB := "zabbix-dc2.yourcompany.com"
	fmt.Println("\n" + "Available zabbix servers are: ")
	fmt.Println(serverA + "\n" + serverB)
	fmt.Println("\n" + "Press 1 for: " + serverA + "\n" + "Press 2 for: " + serverB + "\n")
	fmt.Printf("Enter your option: ")
	reader := bufio.NewReader(os.Stdin)
	optionSelected, _ := reader.ReadString('\n')
	optionSelected = strings.TrimSpace(optionSelected)

	switch optionSelected {
	case "1":
		serverSelected = serverA
	case "2":
		serverSelected = serverB
	default:
		panic("Invalid selection!")
	}

	user, password := GetCreds()
	api := zabbix.NewAPI("https://" + serverSelected + "/api_jsonrpc.php")
	_, err := api.Login(user, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return *api
}

func DeleteHosts(api zabbix.API) {
	a := zabbix.Hosts{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter servers to delete, separated by comma: ")
	hosts, _ := reader.ReadString('\n')
	hosts = strings.TrimSpace(hosts)
	hostList := strings.Split(hosts, ",")
	for _, i := range hostList {
		i = strings.TrimSpace(i)
		b, err := api.HostGetByHost(i)
		if err != nil {
			panic(err)
		}
		a = append(a, *b)
	}
	fmt.Println("Are you Ok to delete below servers?\nName      HostID")
	for _, i := range a {
		fmt.Println(i.Name, i.HostId)
	}
	if Ask4Confirm() {
		api.HostsDelete(a)
	}
}
