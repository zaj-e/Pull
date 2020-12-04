package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var direccion_nodo string

const (
	puerto_registro  = 8000
	puerto_notifica  = 8001
	puerto_proceso   = 8002
	puerto_solicitud = 8003
)

var bitacora = []string{"192.168.101.9", "192.168.101.25"}

type Info struct {
	Tipo     string
	NodeNum  int
	NodeAddr string
}
type MyInfo struct {
	contMsg  int
	first    bool
	nextNum  int
	nextAddr string
}

var puedeIniciar chan bool
var chMiInfo chan MyInfo

var ticket int

func main() {
	//Identificación = IP
	direccion_nodo = descubreIP() //"192.168.101.12"
	fmt.Printf("IP = %s\n", direccion_nodo)
	//Server
	//go registrarServidor()
	//go registrarProceso()
	//Client
	//A que nodo se requiere unir???
	//bufferIn := bufio.NewReader(os.Stdin)
	/*
		fmt.Print("Ingrese el Ip del nodo a unir: ")
		strDirNodoRem, _ := bufferIn.ReadString('\n')
		strDirNodoRem = strings.TrimSpace(strDirNodoRem)
		if strDirNodoRem != "" {
			registrarCliente(strDirNodoRem)
		}
	*/
	//generaciòn de ticket
	rand.Seed(time.Now().UTC().UnixNano())
	ticket = rand.Intn(1000000)
	fmt.Printf("Nro Ticket %d\n", ticket)

	//crear los canales
	puedeIniciar = make(chan bool)
	chMiInfo = make(chan MyInfo)
	//enviar la solicitud inicial
	go func() {
		chMiInfo <- MyInfo{0, true, 1000001, ""}
	}()
	//esperar el inicio de la solicitud
	go func() {
		fmt.Println("Presiona enter para Iniciar la solicitud...")
		bufferIn := bufio.NewReader(os.Stdin)
		msg, _ := bufferIn.ReadString('\n')
		fmt.Println(msg)
		info := Info{"SENDNUM", ticket, direccion_nodo}
		//enviar a todos los nodos de la red
		fmt.Println("paso enviar")
		for _, direccion := range bitacora {
			go enviarSolicitud(direccion, info)
		}
	}()

	escucharSolicitud()
	//servidor
	//escucharNotificaciones()
}

func registrarServidor() {
	//formatear el host para la conexion  ip:port
	hostname := fmt.Sprintf("%s:%d", direccion_nodo, puerto_registro)
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		//aceptar las conexiones
		conn, _ := ln.Accept()
		//manejar la conexion concurrentemente
		go manejadorRegistro(conn)
	}
}
func manejadorRegistro(conn net.Conn) {
	defer conn.Close()
	//recuperar el mensaje de registro (IP de nodo)
	//Leer el Ip del nodo q solicitó unirse
	bufferIn := bufio.NewReader(conn)
	msgIP, _ := bufferIn.ReadString('\n')
	msgIP = strings.TrimSpace(msgIP)
	//serializar la bitacora de direcciones
	bytesBitacora, _ := json.Marshal(bitacora)
	fmt.Fprintf(conn, "%s\n", string(bytesBitacora)) //escribiendo al cliente
	notificarTodos(msgIP)                            //Pull del mensaje hacia todos los nodos registrados en la bitacora
	bitacora = append(bitacora, msgIP)               //Agregar la ip que llega como mensaje, a la bitacora
	fmt.Println(bitacora)                            //imprimir la bitacora
}
func notificarTodos(msgIP string) {
	//hace un pull del mensaje
	for _, direccion := range bitacora {
		notificar(direccion, msgIP)
	}
}
func notificar(direccion, msgIP string) {
	remotehost := fmt.Sprintf("%s:%d", direccion, puerto_notifica)
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", msgIP) //enviar el msgIp por la conexión del cliente
}
func registrarCliente(strDirNodoRem string) {
	remotehost := fmt.Sprintf("%s:%d", strDirNodoRem, puerto_registro)
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	//el cliente envía su dirección IP
	fmt.Fprintf(conn, "%s\n", direccion_nodo)
	//recibe la bitacora del servidor
	bufferIn := bufio.NewReader(conn)
	msgBitacora, _ := bufferIn.ReadString('\n')
	//deserializar
	var arrBitacora []string
	json.Unmarshal([]byte(msgBitacora), &arrBitacora)
	bitacora = append(arrBitacora, strDirNodoRem) //actualiza la bitacora de direcciones ip del cluster
	fmt.Println(bitacora)
}
func escucharNotificaciones() {
	hostname := fmt.Sprintf("%s:%d", direccion_nodo, puerto_notifica)
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go manejadorNotificaciones(conn)
	}
}
func manejadorNotificaciones(conn net.Conn) {
	defer conn.Close()
	//recuperar el mensaje IP
	bufferIn := bufio.NewReader(conn)
	msgIP, _ := bufferIn.ReadString('\n')
	msgIP = strings.TrimSpace(msgIP)
	//agregarlo a la bitacora
	bitacora = append(bitacora, msgIP)
	fmt.Println(bitacora)
}
func descubreIP() string {
	listaInterfaces, _ := net.Interfaces()
	for _, interf := range listaInterfaces {
		//fmt.Println(interf.Name)
		direcciones, _ := interf.Addrs()
		for _, direccion := range direcciones {
			//fmt.Println(direccion)
			switch d := direccion.(type) {
			case *net.IPNet:
				//fmt.Println(d.IP)
				if strings.HasPrefix(d.IP.String(), "192") {
					//fmt.Println(d.IP)
					return d.IP.String()
				}
			}

		}
	}
	return ""
}

func registrarProceso() {
	hostname := fmt.Sprintf("%s:%d", direccion_nodo, puerto_proceso)
	ln, _ := net.Listen("tcp", hostname) //IP:Port
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go manejarProceso(conn)
	}
}

func manejarProceso(conn net.Conn) {
	defer conn.Close()
	bufferIn := bufio.NewReader(conn)
	strNum, _ := bufferIn.ReadString('\n')
	strNum = strings.TrimSpace(strNum)
	num, _ := strconv.Atoi(strNum)
	fmt.Printf("Numero recibido: %d\n", num)
	if num == 0 {
		fmt.Println("Boommm!!!")
	} else {
		enviarNumero(num - 1)
	}
}
func enviarNumero(num int) {
	indice := rand.Intn(len(bitacora))
	remoteHost := fmt.Sprintf("%s:%d", bitacora[indice], puerto_proceso)
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()
	//enviar
	fmt.Printf("\nEnvió %d hacia %s\n", num, bitacora[indice])
	fmt.Fprintf(conn, "%d\n", num)
}

func escucharSolicitud() {
	hostname := fmt.Sprintf("%s:%d", direccion_nodo, puerto_solicitud)
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go manejadorSolicitudes(conn)
	}
}
func manejadorSolicitudes(conn net.Conn) {
	defer conn.Close()
	bufferIn := bufio.NewReader(conn)
	msg, _ := bufferIn.ReadString('\n')
	//leer
	var info Info
	json.Unmarshal([]byte(msg), &info)
	fmt.Println(info)
	switch info.Tipo {
	case "SENDNUM":
		myInfo := <-chMiInfo
		if info.NodeNum < ticket {
			myInfo.first = false
		} else if info.NodeNum < myInfo.nextNum {
			myInfo.nextNum = info.NodeNum
			myInfo.nextAddr = info.NodeAddr
		}
		myInfo.contMsg++
		go func() {
			chMiInfo <- myInfo
		}()
		//autorizacion
		if myInfo.contMsg == len(bitacora) {
			if myInfo.first {
				//atender la seccion critica
				atenderSeccionCritica()
			} else {
				puedeIniciar <- true
			}
		}
	case "START":
		<-puedeIniciar
		//atender seccion critica
		atenderSeccionCritica()
	}
}

func enviarSolicitud(direccion string, msg Info) {
	//envìa la info a todos los nodos
	remoteHost := strings.TrimSpace(direccion)
	remoteHost = fmt.Sprintf("%s:%d", remoteHost, puerto_solicitud)
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()
	bMsg, _ := json.Marshal(msg)
	fmt.Fprintln(conn, string(bMsg))
}
func atenderSeccionCritica() {
	fmt.Println("Iniciando trabajo en seccion critica")
	myInfo := <-chMiInfo
	if myInfo.nextAddr == "" {
		fmt.Println("Soy el proceso unico")
	} else {
		fmt.Println("Finalizando trabajo en la seccion critica")
		fmt.Printf("Siguiente a procesar es %s con el ticket %d\n", myInfo.nextAddr, myInfo.nextNum)
		msg := Info{Tipo: "START"}
		enviarSolicitud(myInfo.nextAddr, msg) // se esta comunicando al nodo su orden de acceso a la SC
	}
}
