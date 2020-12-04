package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type tipoMsg struct {
	Code int
	Addr string
	Op   int
}

const (
	cNum = iota
	opA  = 1
	opB  = 2
)

var direccion_nodo string
var bitacora = []string{"localhost:9000", "localhost:9002"}

var chInfo chan map[string]int //sincronizar procesos

func main() {
	direccion_nodo = "localhost:9001"
	chInfo = make(chan map[string]int)
	go func() { chInfo <- map[string]int{} }()
	go server()
	//Evaluar
	time.Sleep(time.Millisecond * 100)
	var opinion int
	for {
		fmt.Print("\nIngrese opinion: ")
		fmt.Scanf("%d\n", &opinion)
		mensaje := tipoMsg{cNum, direccion_nodo, opinion}
		for _, direccion := range bitacora {
			enviar(direccion, mensaje)
		}
	}
}

func server() {
	//rol de escuchar
	if ln, errores := net.Listen("tcp", direccion_nodo); errores != nil {
		log.Panicln("No se puede iniciar la conectividad", direccion_nodo)
	} else {
		defer ln.Close()
		fmt.Printf("Escuchando en el host %s\n", direccion_nodo)
		for {
			//aceptar las conexiones
			if conn, err := ln.Accept(); err != nil {
				log.Panicln("Error al aceptar conexion", conn.RemoteAddr())
			} else {
				go manejadorMensajes(conn)
			}
		}
	}
}
func manejadorMensajes(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)

	var msg tipoMsg
	if errores := dec.Decode(&msg); errores != nil {
		log.Panicln(" No se pudo decodificar desde ", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case cNum:
			consenso(msg)
		}
	}
}

func consenso(msg tipoMsg) {
	info := <-chInfo
	info[msg.Addr] = msg.Op // opinion
	//criterio de evaluaciòn
	if len(info) == len(bitacora) { //verificar si todos los nodos enviaron opinion
		ca, cb := 0, 0
		for _, op := range info {
			if op == opA {
				ca++
			} else {
				cb++
			}
		}
		//resultado del consenso
		if ca > cb {
			fmt.Println("Repuesta A!")
		} else {
			fmt.Println("Respuesta B!")
		}
		info = map[string]int{}
	}
	go func() {
		chInfo <- info
	}()
}

func enviar(direccion string, msg tipoMsg) {
	if conn, err := net.Dial("tcp", direccion); err != nil {
		log.Panicln("no se puede establecer conexion", direccion)
	} else {
		defer conn.Close()
		fmt.Println("Enviando mensaje a: ", direccion)
		envioMsg := json.NewEncoder(conn)
		envioMsg.Encode(msg)
	}

}
