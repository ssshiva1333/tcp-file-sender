package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	//"io"
)

const (
	buffer = 1024
)

type required_info struct {
	port string
	receiver_ip string
	file_path string
}

func main() {
	fmt.Println("** File sender which uses TCP **")
	fmt.Println("**    Developer -> shiva13    **")
	fmt.Println()
	fmt.Println("* Parameter line -> \"receivr ip\" port \"save path\"")
	fmt.Println()

	//required_info := &required_info{receiver_ip: os.Args[0], port: os.Args[1], file_path: os.Args[2]}

	required_info := &required_info{receiver_ip: "localhost", port: "12345", file_path: "C:\\Users\\user\\Desktop\\metasploit.docx"}

	required_info.send_file_opertions()
}

func (r *required_info) send_file_opertions() {
	tcp_address := r.address_resolver()

	connection := connect_to(tcp_address)
	defer connection.Close()

	file_name := r.file_name_extracter()

	send_file_name(file_name, connection)

	r.read_file_and_send(connection)
}

func (r *required_info) address_resolver() (tcp_address *net.TCPAddr) {
	tcp_address, address_error := net.ResolveTCPAddr("tcp", r.receiver_ip + ":" + r.port)
	if address_error != nil {
		fmt.Println(address_error.Error())

		os.Exit(1)
	}

	return
}

func connect_to(tcp_address *net.TCPAddr) (*net.TCPConn) {
	var connection *net.TCPConn
	var connection_error error

	for {
		connection, connection_error = net.DialTCP("tcp", nil, tcp_address)			
		if connection_error != nil {
			fmt.Println(connection_error.Error())

			continue
		} else {
			break
		}
	}

	return connection
}

func (r *required_info) read_file_and_send(connection *net.TCPConn) {
	file, opening_error := os.Open(r.file_path)
	if opening_error != nil {
		fmt.Println(opening_error.Error())

		os.Exit(1)
	}
	defer file.Close()

	bytes := make([]byte, 1024)
	for {
		n, reading_error := file.Read(bytes)
		if reading_error != nil {
			fmt.Println(reading_error.Error())
		}

		if n == 0 {
			break
		}

		connection.Write(bytes[:n])		
	}
}

func send_file_name(file_name string, connection *net.TCPConn) {
	_, writing_error := connection.Write([]byte(file_name))
	if writing_error != nil {
		fmt.Println(writing_error.Error())

		os.Exit(1)
	}
}

func (r *required_info) file_name_extracter() (name string) {
	file_name_index := strings.LastIndex(r.file_path, "/")
	if file_name_index != -1 {
		name = r.file_path[file_name_index + 1 :len(r.file_path)]
	} else {
		name = r.file_path[strings.LastIndex(r.file_path, "\\") + 1 :len(r.file_path)]
	}

	return
}