package main

import (
	"fmt"
	"net"
)

func main() {
	//establecemos la conexion al nodo master
	con, _ := net.Dial("tcp", "localhost:8000")
	defer con.Close()
	//Que enviamos???
	fmt.Fprintln(con, "Comunicando desde el nodo cliente!!!")
}
