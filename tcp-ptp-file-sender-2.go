package main

import (
	"fmt"
	"net"
	"os"
	"io"
)

const (
	buffer = 1024
)

type required_info struct {
	port string
	save_path string
}

func main() {
	fmt.Println("** File sender which uses TCP **")
	fmt.Println("**    Developer -> shiva13    **")
	fmt.Println()
	fmt.Println("* Parameter line -> port \"save path\"")
	fmt.Println()

	//required_info := &required_info{port: os.Args[0], save_path: os.Args[1]}

	required_info := &required_info{port: "12345", save_path: "C:\\Users\\user\\Desktop\\"}

	required_info.taking_file_operations()
	
}

func (r *required_info) taking_file_operations() {
	tcp_address := r.address_resolver()

	listener := listener(tcp_address)
	defer listener.Close()
	
	for {
		connection := accept_request(listener)

		r.take_file_and_write(connection)
	}
}

func (r *required_info) address_resolver() (tcp_address *net.TCPAddr) {
	tcp_address, address_error := net.ResolveTCPAddr("tcp", "0.0.0.0:" + r.port)
	if address_error != nil {
		fmt.Println(address_error.Error())
		os.Exit(1)
	}

	return
}

func listener(tcp_address *net.TCPAddr) (*net.TCPListener) {
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

func (r *required_info) take_file_and_write(connection *net.TCPConn) {
	file_name := ""
	var file *os.File
	var creating_error error

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

	fmt.Println("* It has been done")
}