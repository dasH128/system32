package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var nodos map[string]bool
var host string
var puerto int
var remotehost string

func main() {
	//fmt.Print("Dirección local [host]: ")
	host = "localhost" //fmt.Scanln(&host)

	fmt.Print("Dirección local [puerto]: ")
	fmt.Scanln(&puerto)

	var remote string
	fmt.Print("conetarse a [puerto]: ")
	fmt.Scanln(&remote)

	nodos = make(map[string]bool)

	go ServidorAgregador()

	/*for {
		log.Println("Test")
	}*/
	// nodos["localhost:9000"] = true
	// nodos["localhost:9001"] = true
	// nodos["localhost:9002"] = true
	// ClienteAgregador("localhost:8000")
	//go EnviarSinRespuesta("test", "target")

}
func EnviarSinRespuesta(msj string, target string) {

	clienteHttp := &http.Client{}
	hostServer := host + ":" + strconv.Itoa(puerto+1)

	usuarioComoJson, _ := json.Marshal(msj)
	// usuarioComoJson := "Pruebaaaaaaaaaaaaa"
	peticion, _ := http.NewRequest("POST", hostServer, bytes.NewBuffer(usuarioComoJson))
	respuesta, _ := clienteHttp.Do(peticion)
	defer respuesta.Body.Close()

	cuerpoRespuesta, _ := ioutil.ReadAll(respuesta.Body)
	respuestaString := string(cuerpoRespuesta)
	log.Printf("Cuerpo de respuesta del servidor: '%s'", respuestaString)
}
func EnviarConRespuesta(msj string, target string) {

	clienteHttp := &http.Client{}
	hostServer := host + ":" + strconv.Itoa(puerto+1)

	usuarioComoJson, _ := json.Marshal(msj)
	// usuarioComoJson := "Pruebaaaaaaaaaaaaa"
	peticion, _ := http.NewRequest("POST", hostServer, bytes.NewBuffer(usuarioComoJson))
	respuesta, _ := clienteHttp.Do(peticion)
	defer respuesta.Body.Close()

	cuerpoRespuesta, _ := ioutil.ReadAll(respuesta.Body)
	respuestaString := string(cuerpoRespuesta)
	log.Printf("Cuerpo de respuesta del servidor: '%s'", respuestaString)
}
func AgregarHost(host string) {
	nodos[host] = true
}
func ServidorAgregador() {
	hostServer := host + ":" + strconv.Itoa(puerto+1)
	http.HandleFunc("/", HandleAgregador)

	if err := http.ListenAndServe(hostServer, nil); err != nil {
		log.Fatal(err)
	}
}
func HandleAgregador(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var hostname string
		_ = json.NewDecoder(r.Body).Decode(&hostname)

		for nodo, _ := range nodos {
			AvisarNodo(nodo, hostname)
		}

		dataJson, _ := json.Marshal(hostname)
		w.Header().Set("Content-Type", "application/json")
		AgregarHost(hostname)
		w.Write(dataJson)
		fmt.Print("Estas enviando ", hostname)
	}
	//fmt.Println("handle")
}
func AvisarNodo(nodo string, nuevo string) {

}
func ClienteAgregador(nuevoHost string) {
	for k, v := range nodos {
		fmt.Println(k, "=>", v)
		EnviarSinRespuesta(nuevoHost, k)
	}

}
func ServidorRegistrador() {
	hostServer := host + ":" + strconv.Itoa(puerto+2)
	http.HandleFunc("/", HandleRegistrador)

	if err := http.ListenAndServe(hostServer, nil); err != nil {
		log.Fatal(err)
	}
}
func HandleRegistrador(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data string
		_ = json.NewDecoder(r.Body).Decode(&data)
		//agregar el nuevo nodo a su arreglo de

	}
}
func ClienteRegistrador() {

}
