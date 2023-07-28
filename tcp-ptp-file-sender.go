package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	buffer = 1024
)

// STRUCT FOR SENDING FILE
type required_info_SEND struct {
	port        string
	receiver_ip string
	file_path   string
}

// STRUCT FOR TAKING FILE
type required_info_TAKE struct {
	port      string
	save_path string
}

func main() {
	fmt.Println("** File sender which uses TCP **")
	fmt.Println("**    Developer -> shiva13    **")
	fmt.Println()

	fmt.Print("* Port number -> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	port := scanner.Text()

	for {
		fmt.Println("* Press \"1\" to send file")
		fmt.Println("* Press \"2\" to take file")
		fmt.Println()

		fmt.Print("* > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		fmt.Println()
		fmt.Println()

		if scanner.Text() == "1" {
			var receiver_ip string
			var file_path string

			fmt.Println("* Parameter line -> \"receiver ip\" \"file path\"")
			fmt.Print("* > ")

			_, err := fmt.Scanln(&receiver_ip, &file_path)
			if err != nil {
				fmt.Println("* ", err.Error())

				break
			}

			required_info_SEND := &required_info_SEND{receiver_ip: receiver_ip, port: port, file_path: file_path}
			required_info_SEND.send_file_opertions()
		} else if scanner.Text() == "2" {
			var save_path string

			fmt.Println("* Parameter line ->\"save path\"")
			fmt.Print("* > ")

			_, err := fmt.Scanln(&save_path)
			if err != nil {
				fmt.Println("* ", err.Error())

				break
			}

			required_info_TAKE := &required_info_TAKE{port: port, save_path: save_path}
			required_info_TAKE.taking_file_operations()
		} else {
			fmt.Println("* Wrong option")
		}
	}
}

//FILE SENDING FUNCTIONS

func (r *required_info_SEND) send_file_opertions() {
	tcp_address := r.address_resolver()

	connection := connect_to(tcp_address)
	defer connection.Close()

	file_name := r.file_name_extracter()

	send_file_name(file_name, connection)

	r.read_file_and_send(connection)
}

func (r *required_info_SEND) address_resolver() (tcp_address *net.TCPAddr) {
	tcp_address, address_error := net.ResolveTCPAddr("tcp", r.receiver_ip+":"+r.port)
	if address_error != nil {
		fmt.Println(address_error.Error())

		os.Exit(1)
	}

	return
}

func connect_to(tcp_address *net.TCPAddr) *net.TCPConn {
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

func (r *required_info_SEND) read_file_and_send(connection *net.TCPConn) {
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

func (r *required_info_SEND) file_name_extracter() (name string) {
	file_name_index := strings.LastIndex(r.file_path, "/")
	if file_name_index != -1 {
		name = r.file_path[file_name_index+1 : len(r.file_path)]
	} else {
		name = r.file_path[strings.LastIndex(r.file_path, "\\")+1 : len(r.file_path)]
	}

	return
}

//FILE TAKING FUNCTIONS

func (r *required_info_TAKE) taking_file_operations() {
	tcp_address := r.address_resolver()

	listener := listener(tcp_address)
	defer listener.Close()

	for {
		connection := accept_request(listener)

		n := r.take_file_and_write(connection)
		if n == 1 {
			break
		}
	}
}

func (r *required_info_TAKE) address_resolver() (tcp_address *net.TCPAddr) {
	tcp_address, address_error := net.ResolveTCPAddr("tcp", "0.0.0.0:"+r.port)
	if address_error != nil {
		fmt.Println(address_error.Error())

		os.Exit(1)
	}

	return
}

func listener(tcp_address *net.TCPAddr) *net.TCPListener {
	var listener *net.TCPListener
	var listener_error error

	for {
		listener, listener_error = net.ListenTCP("tcp", tcp_address)
		if listener_error != nil {
			fmt.Println(listener_error.Error())

			continue
		} else {
			break
		}
	}

	return listener
}

func accept_request(listener *net.TCPListener) (connection *net.TCPConn) {
	connection, connection_error := listener.AcceptTCP()
	if connection_error != nil {
		fmt.Println(connection_error.Error())

		os.Exit(1)
	}

	return
}

func (r *required_info_TAKE) take_file_and_write(connection *net.TCPConn) int {
	file_name := ""
	var file *os.File
	var creating_error error
	defer file.Close()

	i := 0
	for {
		bytes := make([]byte, 1024)
		n, reading_error := connection.Read(bytes)
		if reading_error != nil {
			fmt.Println(reading_error.Error())
		}

		if reading_error == io.EOF {
			break
		}

		if i == 0 {
			file_name = "sent_file_ " + string(bytes[:n])
			bytes = nil

			file, creating_error = os.Create(r.save_path + file_name)
			if creating_error != nil {
				fmt.Println(creating_error.Error())

				os.Exit(1)
			}

			i++
		}

		_, writing_error := file.Write(bytes)
		if writing_error != nil {
			fmt.Println(writing_error.Error())

			os.Exit(1)
		}
	}

	return 1
}
