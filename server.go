package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
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

func handle(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	fmt.Println("Conexión establecida, esperando su instrucción...")
	msg, _ := r.ReadString('\n')
	fmt.Println("Recibido: ", msg)
	routes := []Route{Route{IP: "123456", Port: "8001"}, Route{IP: "654321", Port: "8001"}}
	node, _ := json.Marshal(routes)
	fmt.Fprintln(con, string(node))
	fmt.Println("Respuesta enviada.")
	time.Sleep(time.Second)
	os.Exit(0)
}

func main() {
	fmt.Println("Iniciando Servidor...")
	ln, _ := net.Listen("tcp", "192.168.1.14:8001")
	defer ln.Close()
	fmt.Println("Escuchando por puerto 8001")
	for {
		con, _ := ln.Accept()
		go handle(con)
	}
}
