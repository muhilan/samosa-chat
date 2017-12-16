package main

import ( "net/http"
   "golang.org/x/net/http2"
   "fmt"
   "log"
   )

func main(){

	var srv http.Server

	srv.Addr = ":8080"

	http2.ConfigureServer(&srv, nil)
	http.HandleFunc("/", handlerHtml)
	fmt.Println("Server started on 8080 ")
	err := srv.ListenAndServeTLS("certs/localhost.cert", "certs/localhost.key")
	log.Fatal(err)

}


func handlerHtml(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1><center> Hello from Go! </h1></center>"))
}