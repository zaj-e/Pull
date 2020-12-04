package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

var direccionNodo string
var direcciones []string //registro de direcciones de los nodos del cluster
const (
	port_registro = 8000 //solicitar el registro de un nuevo nodo al cluster
	port_notifica = 8001 //notificar a toda la red del cluster la incorporación de un nuevo nodo
)

func main() {
	//Indentificacion como servidor
	direccionNodo = myIp() //"192.168.101.12" //direccion local
	fmt.Println(direccionNodo)
	//rol servidor
	go registrarServidor()
	//rol cliente
	//solicita unirse a uno de los nodos existentes en el cluster
	bIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese la direccion remota: ")
	direccionUnir, _ := bIn.ReadString('\n')
	direccionUnir = strings.TrimSpace(direccionUnir)
	//si no se va a unir a ningun nodo entonces es un server
	if direccionUnir != "" {
		registrarCliente(direccionUnir)
	}
	//componrtamiento de escucha constante de notificaciones
	EscucharNotificaciones()
}

func registrarServidor() {
	hostname := fmt.Sprintf("%s:%d", direccionNodo, port_registro)
	//escuchar
	ln, _ := net.Listen("tcp", hostname) //hostname=  direccionIp:port
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go manejadorRegistro(conn)
	}
}

func manejadorRegistro(conn net.Conn) {
	defer conn.Close()
	bIn := bufio.NewReader(conn)
	direccionIP, _ := bIn.ReadString('\n')
	direccionIP = strings.TrimSpace(direccionIP)
	//transmitir la respuesta de la lista de direcciones q guarda el nodo
	//hacia el nuevo nodo
	//serializar la bitacora actual
	bytesDirecciones, _ := json.Marshal(direcciones)
	fmt.Fprintf(conn, "%s\n", string(bytesDirecciones)) //al nuevo nodo
	//comunicar al resto de nodos de la red
	notificarTodos(direccionIP)
	direcciones = append(direcciones, direccionIP)
	fmt.Println(direcciones)
}
func notificarTodos(direccionIP string) {
	for _, dirRemote := range direcciones {
		notificar(dirRemote, direccionIP) //PULL
	}
}
func notificar(dirRemote, direccionIP string) {
	remoteHost := fmt.Sprintf("%s:%d", dirRemote, port_notifica)
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", direccionIP) //la ip del nuevo nodo
}

func registrarCliente(direccioUnir string) {
	remotehost := fmt.Sprintf("%s:%d", direccioUnir, port_registro)
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	//enviar la direccion ip
	fmt.Fprintf(conn, "%s\n", direccionNodo)
	//espera la respuesta del servidor
	//llega la lista de direcciones
	bIn := bufio.NewReader(conn)
	msgListaDir, _ := bIn.ReadString('\n')
	//deserializar, nos envia en formato json
	var respDir []string
	json.Unmarshal([]byte(msgListaDir), &respDir)
	direcciones = append(direcciones, direccioUnir) //agregar la direccion del nodo servidor
	fmt.Println(direcciones)
}

func EscucharNotificaciones() {
	hostname := fmt.Sprintf("%s:%d", direccionNodo, port_notifica)
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go manejadorNotificaciones(conn)
	}
}

func manejadorNotificaciones(conn net.Conn) {
	defer conn.Close()
	bIn := bufio.NewReader(conn)
	dirIP, _ := bIn.ReadString('\n')
	dirIP = strings.TrimSpace(dirIP)
	//Agregar la Ip a la lista de direcciones
	direcciones = append(direcciones, dirIP)
	fmt.Println(direcciones)
}

func myIp() string {
	interfacess, _ := net.Interfaces()
	for _, interf := range interfacess {
		direcciones, _ := interf.Addrs()
		for _, direccion := range direcciones {
			//fmt.Println(direccion)
			switch d := direccion.(type) {
			case *net.IPNet:
				if strings.HasPrefix(d.IP.String(), "192") {
					//fmt.Println(d.IP)
					return d.IP.String()
				}
			}
		}
	}
	return ""
}
