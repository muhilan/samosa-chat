package main

import ( "net/http"
   // "golang.org/x/net/http2"
   "fmt"
   "log"
    "encoding/json"
    "io/ioutil"
   )

type Message struct {
	Owner   string  `json:"Owner"`
	Time string `json:"Time"`
	Text string `json:"Text"`
}


func main(){

	var srv http.Server

	srv.Addr = ":8080"

	// http2.ConfigureServer(&srv, nil)
	http.HandleFunc("/", handler)
	fmt.Println("Server started on 8080 ")
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

