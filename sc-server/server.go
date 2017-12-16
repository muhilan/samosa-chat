package main

import ( "net/http"
 
   "golang.org/x/net/http2"

   "fmt"

   )

func main(){

	fmt.Println("Server")

	var srv http.Server

	srv.Addr = ":8081"

	

	//Enable http2

	http2.ConfigureServer(&srv, nil)
		http.HandleFunc("/", handlerHtml)
	fmt.Println("Server started on 8081 ")
	srv.ListenAndServeTLS("certs/localhost.cert", "certs/localhost.key")

	




}


func handlerHtml(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1><center> Hello from Go! </h1></center>"))
}