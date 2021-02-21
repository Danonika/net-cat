package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	st       = map[string]net.Conn{}
	messages string
	mutex    sync.Mutex
)

/*__________________________________________________*/
//Get our time
func getTime() string {
	date := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
}

/*__________________________________________________*/
//Sent messages to all TCP connections
func sent(message, name string) {
	for i := range st {
		if i != name {
			st[i].Write([]byte("\n" + message))
		}
		st[i].Write([]byte(fmt.Sprintf("[%s][%s]:", getTime(), i)))
	}
	messages += message
}

/*__________________________________________________*/
//Handle user
func chat(user net.Conn) {
	if len(st) == 10 {
		user.Write([]byte("Sorry chat is full. Try again later"))
		user.Close()
		return
	}
	var name, message string
	f, _ := os.Open("logo.txt")
	r := bufio.NewReader(f)
	buffer := make([]byte, 500)
	io.ReadFull(r, buffer)
	// welcome, _ := ioutil.ReadFile("logo.txt")
	// io.ReadFull()
	user.Write(buffer)
	// user.Write(welcome)
	user.Write([]byte("[ENTER YOUR NAME]:"))
	scanner := bufio.NewScanner(user)
	/*__________________________________________________*/

	//User name
	for scanner.Scan() {
		name = scanner.Text()
		name = strings.TrimSpace(name)
		_, ok := st[name]
		if len(name) == 0 || ok {
			user.Write([]byte("Please enter correct name:"))
		} else {
			st[name] = user
			break
		}
	}

	/*__________________________________________________*/
	//Chat handler
	user.Write([]byte(messages))
	sent(fmt.Sprintf("%s has joined out chat ...\n", name), name)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		text := scanner.Text()
		check := strings.Trim(text, " \n\t\r")
		if len(check) != 0 {
			mutex.Lock()
			message = string(fmt.Sprintf("[%s][%s]:", getTime(), name)) + text + "\n"
			sent(message, name)
			mutex.Unlock()
		}
	}
	//Ops user left
	sent(fmt.Sprintf("%s has left out chat ...\n", name), name)
	delete(st, name)
}
func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	} else {
		if len(args) == 0 {
			args = append(args, "8989")
		}
		fmt.Printf("Listening on the port :%s", args[0])
		ln, _ := net.Listen("tcp", ":"+args[0])
		for {
			connection, _ := ln.Accept()
			go chat(connection)
		}
	}

}
