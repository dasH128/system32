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

type HostName struct {
	Host string
	Port int
}

var nodos []string
var hostname HostName

func main() {
	hostname.Host = "loclahost"
	fmt.Print("Dirección local [puerto]: ")
	fmt.Scanln(&hostname.Port)

	// var remote string
	// fmt.Print("conetarse a [host:puerto]: ")
	// fmt.Scanln(&remote)

	nodos = append(nodos, "localhost:3001")
	nodos = append(nodos, "localhost:3002")
	ServerRegistro()

}

func ServerRegistro() {
	//hostServer := hostname.Host + ":" + strconv.Itoa((hostname.Port)+1)
	http.HandleFunc("/", HandleRegistro)
	hostServer := hostname.Host + ":" + strconv.Itoa(hostname.Port+1)
	//var h = hostname.Host + ":" + hostname.Port
	fmt.Println(hostServer)
	if err := http.ListenAndServe(hostServer, nil); err != nil {
		log.Fatal(err)
	}
}

func HandleRegistro(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data string
		_ = json.NewDecoder(r.Body).Decode(&data)

		//avisar a los demas nodos
		for nodo, _ := range nodos {
			fmt.Println(nodos[nodo])
			AvisarNodo(nodos[nodo], data)
		}

		//enviar la lista mas mi hpost
		dataJson, _ := json.Marshal(nodos)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataJson)
		fmt.Print("Los host en la red son: ", nodos)

		//agregar su host a mi lista
		nodos = append(nodos, data)
		fmt.Println(nodos)
	}
}
func AvisarNodo(nodo string, hostname string) {
	clienteHttp := &http.Client{}

	dataComoJson, _ := json.Marshal(nodo)
	peticion, _ := http.NewRequest("POST", hostname, bytes.NewBuffer(dataComoJson))
	respuesta, _ := clienteHttp.Do(peticion)
	defer respuesta.Body.Close()

	cuerpoRespuesta, _ := ioutil.ReadAll(respuesta.Body)
	respuestaString := string(cuerpoRespuesta)
	log.Printf("Cuerpo de respuesta del servidor: '%s'", respuestaString)
}

func ServerNotificaciones() {
	http.HandleFunc("/", HandleNotificaciones)
	hostServer := hostname.Host + ":" + strconv.Itoa(hostname.Port+2)
	fmt.Println(hostServer)
	if err := http.ListenAndServe(hostServer, nil); err != nil {
		log.Fatal(err)
	}
}

func HandleNotificaciones(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var nuevoNodo string
		_ = json.NewDecoder(r.Body).Decode(&nuevoNodo)
	}
}
