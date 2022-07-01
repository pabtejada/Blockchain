package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Route struct {
	IP   string `json: "IP"`
	Port string `json:"Port"`
}

type Node struct {
	Route Route `json: "Route"`
	// 0 = nuevo, 1 = solo agrega
	Instruction int `json:"Instruction"`
}

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func getPrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	check(err)
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	log.Fatal("No hay ip disponible")
	os.Exit(1)
	return ""
}

func connectToNetwork(sourceRoute Route, destinationRoute Route) {
	log.Println("Uniendose a la red...")
	con, err := net.Dial("tcp", fmt.Sprintf("%s:%s", destinationRoute.IP, destinationRoute.Port))
	check(err)
	defer con.Close()
	node, err := json.Marshal(Node{Route: sourceRoute, Instruction: 0})
	check(err)
	//log.Println(string(node))
	log.Println("Enviando su IP...")
	fmt.Fprintln(con, string(node))
	log.Println("Recibiendo las otras IPs...")
	r := bufio.NewReader(con)
	//time.Sleep(500 * time.Millisecond) // comentar solo para pruebas
	resp, err := r.ReadString('\n')
	check(err)
	//log.Println("Data:", resp)
	json.Unmarshal([]byte(resp), &routes)
	//log.Println(routes)
	log.Println("Completado.")
}

const PORT string = "8001"

var routes = make([]Route, 0)

func main() {
	// Pasar la ip de algun nodo por argumento del programa
	// ej: 'go run conectarse_a_la_red.go 192.168.1.50'
	destinationIP := os.Args[1]
	myRoute := Route{IP: getPrivateIP(), Port: PORT}
	log.Println("NODE IP: ", myRoute.IP)
	destinationRoute := Route{IP: destinationIP, Port: PORT}
	connectToNetwork(myRoute, destinationRoute)
}
