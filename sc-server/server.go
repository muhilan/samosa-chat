package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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

	var srv http.Server

	srv.Addr = ":8080"

	// http2.ConfigureServer(&srv, nil)
	http.HandleFunc("/", hiJack)
	//http.HandleFunc("/hijack", hiJack)
	fmt.Println("Server started on 8080 ")
	go func() {
		for {
			select {
			case connCtx := <-activeConns:
				fmt.Printf("Handling connection for %s ", connCtx.owner)
				_ , err := connCtx.connection.Write([]byte("Output"))
				if err != nil {
					deadConns <- connCtx.connection
				}
			case message := <-messages:
				fmt.Printf("Received Msg : %s \n", message)
				for conn, clientId := range connMap {
					go func(conn net.Conn, msg string){
						fmt.Printf("Sending message to client : %s , Message => %s \n", clientId, message.Text)
						_,err := conn.Write([]byte(msg))
						if err != nil {
							deadConns <- conn
						}
					}(conn, message.Text)
				}
			case conn := <- deadConns:
				conn.Close()
				delete(connMap, conn)
			}

		}
	}()

	err := srv.ListenAndServeTLS("certs/localhost.cert", "certs/localhost.key")
	log.Fatal(err)



}

func handler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}





func hiJack(w http.ResponseWriter, r *http.Request){
    fmt.Println("Message received")
	b, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "web server doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		deadConns <- conn
		return
	}
	// Don't forget to close the connection:
	//defer conn.Close()
	bufrw.WriteString("")
	bufrw.Flush()

	//connMap[msg.Owner] = conn
	activeConns <- ConnectionContext{connection: conn, owner: msg.Owner}
	connMap[conn] = msg.Owner
	fmt.Println(msg)
	messages <- msg
}