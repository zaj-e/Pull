package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	con, _ := net.Dial("tcp", "localhost:8000")
	defer con.Close()

	//ingreso de datos
	datosIn := bufio.NewReader(os.Stdin)
	r := bufio.NewReader(con)

	for {
		fmt.Print("Ingrese un mensaje: ")
		msg, _ := datosIn.ReadString('\n')

		//enviar mensaje al nodo master
		fmt.Fprint(con, msg)

		//recibir mensaje del nodo master
		resp, _ := r.ReadString('\n')
		fmt.Printf("Respuesta del nodo master: %s", resp)
	}

}
