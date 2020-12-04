package main

import (
	"encoding/json"
	"fmt"
)

type Alumno struct {
	Codigo   string  `json:"code"`
	Nombre   string  `json:"name"`
	Promedio float32 `json:"grade"`
}

func main() {
	alumnos := []Alumno{
		{"2020152", "Juan Jimenez", 17.50},
		{"2020142", "Carlos Mendoza", 15.50},
		{"2020178", "Luis Tapia", 17.45}}

	jsonBytes, _ := json.MarshalIndent(alumnos, "", " ")
	jsonString := string(jsonBytes)

	fmt.Println(jsonString)

	//deserializar
	var listaAlumnos []Alumno
	json.Unmarshal(jsonBytes, &listaAlumnos)
	fmt.Println(listaAlumnos)

}
