package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var remotehost string
var chOrden chan int
var min, n int

func main() {
	//Estblecer el puerto de escucha
	rIng1 := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto de escucha: ")
	port, _ := rIng1.ReadString('\n')
	port = strings.TrimSpace(port)
	hostname := fmt.Sprintf("localhost:%s", port)

	fmt.Print("Ingrese el puerto remoto: ")
	port, _ = rIng1.ReadString('\n')
	port = strings.TrimSpace(port)
	remotehost = fmt.Sprintf("localhost:%s", port)

	fmt.Print("N: ") //tope
	nelementos, _ := rIng1.ReadString('\n')
	nelementos = strings.TrimSpace(nelementos)
	n, _ = strconv.Atoi(nelementos)
	//sincronizacion
	chOrden = make(chan int, 1)
	chOrden <- 0

	//escuchar
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		//manejar llamadas concurrentes
		go manejador(conn)
	}

}

func manejador(conn net.Conn) {
	defer conn.Close()
	//recuperar lo q se está enviando
	r := bufio.NewReader(conn)
	str, _ := r.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(str))

	fmt.Printf("llegó el numero %d\n", num)
	//logica de ordenamiento
	//sincroniza
	orden := <-chOrden
	if orden == 0 {
		min = num
	} else if num < min {
		//enviar los numeros que no cumplen al resto de aplicaciones
		envio(min)
		min = num
	} else {
		envio(num)
	}
	//continuar
	orden++
	if orden == n {
		fmt.Printf("Mostrar el numero Final: %d\n", min)
		orden = 0
	}
	chOrden <- orden
}

func envio(num int) {
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	//envio
	fmt.Fprintf(conn, "%d\n", num)
}
