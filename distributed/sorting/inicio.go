package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var remote string

func main() {
	//solo envia
	//solicitar el puerto de la app a enviar los numeros
	ring := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto de destino: ")
	port, _ := ring.ReadString('\n')
	port = strings.TrimSpace(port)
	remote = fmt.Sprintf("localhost:%s", port)
	enviar(20)
	enviar(10)
	enviar(50)
	enviar(30)
	enviar(40)
}
func enviar(num int) {
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	fmt.Fprintf(conn, "%d\n", num)
}
