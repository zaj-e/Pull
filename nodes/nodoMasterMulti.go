package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	//escucha
	ln, _ := net.Listen("tcp", "localhost:8000")
	defer ln.Close()
	//escucha constante
	for {
		con, _ := ln.Accept()
		go manejador(con) //permite atender concurrentemente las solicitudes de los clientes
	}
}

func manejador(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	for {
		//recuperando el mensaje enviado desde la APP cliente
		msg, _ := r.ReadString('\n')
		fmt.Printf("Recibido: %s", msg)

		//enviando el mensaje de respuesta hacia la APPcliente
		fmt.Fprintf(con, "Conforme!!!\n")
	}
}
