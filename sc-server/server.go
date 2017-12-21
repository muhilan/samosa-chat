package main

import (
	//"net/http"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"crypto/tls"
	"crypto/rand"
	"bufio"
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

	cert, err := tls.LoadX509KeyPair("certs/localhost.cert", "certs/localhost.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	port := "8080"
	service := fmt.Sprintf("0.0.0.0:%s", port)
	listener, err := net.Listen("tcp", service)

	//listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Printf("server: listening on port %s",port)

	go func() {
		for {
			select {
			case connCtx := <-activeConns:
				fmt.Println("Handling connection for  "+ connCtx.owner)
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
								//http.Error(w, err.Error(), 500)
								//return
							}
						}
						messages <- msg
					}
				}(connCtx.connection, connCtx.owner)

			case singleMessage := <-messages:
				fmt.Println("Received Msg : "+ singleMessage.Text)
				for conn, clientId := range connMap {
					go func(conn net.Conn, msg string) {
						fmt.Printf("Sending message to client : %s , Message => %s \n", clientId, singleMessage.Text)
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
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr()) //}
		activeConns <- ConnectionContext{connection: conn}

	}

}

func getJSONString(msgCtx Message) string {
	b, err := json.Marshal(msgCtx)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("Get json : " + string(b) + "\n")
	return string(b) + "\n"
}
