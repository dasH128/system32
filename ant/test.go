package main

import (
	"encoding/json"
	"fmt"
)

type data struct {
	Nombre string `json:"nombre"`
}

func main() {
	var mapas map[string]int

	mapas = make(map[string]int)
	fmt.Println(mapas)

	mapas["uno"] = 1
	fmt.Println(mapas)

	data1 := data{Nombre: "Jordy"}
	fmt.Println(data1)

	mapB, _ := json.Marshal(data1)
	fmt.Println(string(mapB))

	mapC, _ := json.Marshal(mapas)
	fmt.Println(string(mapC))

	num := mapas["uno"]
	fmt.Println(num)
	// dec := json.NewDecoder(&data1)
	// fmt.Println(dec)
}
