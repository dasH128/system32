package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var nodos map[string]bool
var host string
var puerto int

func AgregarHostName(host string) {
	nodos[host] = true
}

func EnviarSinRespuesta(msj string, target string) {
}

func EnviarConRespuesta(msj string, target string) {
}

func Handle(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	strNum, _ := r.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(strNum))
	fmt.Printf("Recibimos el %d\n", num)

}

func ServidorAgregador() {
	hostServer := host + ":" + strconv.Itoa(puerto+1)

	if ln, err := net.Listen("tcp", hostServer); err != nil {
		log.Panicln(err.Error())
	} else {
		defer ln.Close()
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Panicln(err.Error())
			} else {
				go Handle(conn)
			}
		}
	}
}

func EnviarData(num int) {

	fmt.Print("Enviar data")
	hostServer := host + ":" + strconv.Itoa(puerto+1)
	fmt.Print("Enviar data a " + hostServer)
	conn, _ := net.Dial("tcp", hostServer)
	defer conn.Close()
	fmt.Fprintln(conn, num)
}
func main() {
	//fmt.Print("Dirección local [host]: ")
	host = "localhost" //fmt.Scanln(&host)

	//fmt.Print("Dirección local [puerto]: ")
	puerto = 5003 //fmt.Scanln(&puerto)

	//go EnviarData(20)
	ServidorAgregador()

	func() {
		EnviarData(30)
	}()

	//go EnviarData(30)

	//fmt.Scanln(&puerto)

}
