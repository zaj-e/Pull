package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	//escuchar por el puerto 8000
	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Falla al momento de escuchar usando el puerto 8000: ", err.Error())
		os.Exit(1)
	}
	defer ln.Close()
	con, err := ln.Accept()
	if err != nil {
		fmt.Println("Fallo al aceptar la conexion: ", err.Error())
		//continuar con el manejo de esa error
	}
	defer con.Close()
	//lectura de datos que envio el cliente
	r := bufio.NewReader(con)
	msg, _ := r.ReadString('\n')
	fmt.Println(msg)
}
