package main

import (
	"encoding/json"
	"fmt"
	"net"
	"bufio"
	"os"
)

type Message struct {
	Owner string
	Time  int64
	Text  string
}

type ConnectionContext struct {
	connection net.Conn
	owner string
}
var activeConns = make(chan ConnectionContext)
var deadConns = make(chan net.Conn,10)
var connMap = make(map[net.Conn]string)
var messages = make(chan Message)


func main() {
	port :=  "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	service := fmt.Sprintf("0.0.0.0:%s", port)
	listener, err := net.Listen("tcp", service)

	if err != nil {
		fmt.Errorf("server: listen: %s", err)
	}
	fmt.Printf("server: listening on port %s",port)

	go func() {
		for {
			select {
			case connCtx := <-activeConns:
				fmt.Println("Handling connection for  "+ connCtx.owner)
				// There wouldn't be any owner for the first time message
				if connCtx.owner == "" {
					connCtx.owner = "master"
				}
				connMap[connCtx.connection] = connCtx.owner
				go func(conn net.Conn, owner string) {
					reader := bufio.NewReader(conn)
					for {
						fmt.Println("Inside for loop waiting for messages")
						incoming, err := reader.ReadString('\n')
						fmt.Println("Got new msg "+ incoming)
						if err != nil {
							fmt.Println(err.Error())
							deadConns <- connCtx.connection
							break
						}
						var msg Message
						if incoming != "" {
							err = json.Unmarshal([]byte(incoming), &msg)
							if err != nil {
								fmt.Println(err.Error())
							}
						}
						messages <- msg
					}
				}(connCtx.connection, connCtx.owner)

			case singleMessage := <-messages:
				for conn, clientId := range connMap {
					go func(conn net.Conn, msg string) {
						fmt.Printf("Sending message to client : %s \n", clientId)
						_, err := conn.Write([]byte(getJSONString(singleMessage)))
						if err != nil {
							deadConns <- conn
						}
					}(conn, singleMessage.Text)
				}
			case conn := <-deadConns:
				conn.Close()
				delete(connMap, conn)
			}

		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("server: accept: %s", err)
			break
		}
		fmt.Printf("server: accepted from %s", conn.RemoteAddr()) //}
		activeConns <- ConnectionContext{connection: conn}

	}

}

func getJSONString(msgCtx Message) string {
	b, err := json.Marshal(msgCtx)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b) + "\n"
}
