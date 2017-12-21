# WIP : samosa-chat

An extremely simple chat Application in Go


## Server 
 - runs on 8080 by default

```
cd sc-server
go get
go run server.go
````

## Client
#### Create a file '.samosa-chat.json' under your home with the following content

````
{
 "Owner" : "Your Name",
 "OwnerEmail" : "your email",
 "ChatServerHost" : "Server host",
 "ChatServerPort" : "Server port"
}
````
and to start your client app 
````
cd sc-client
go get
go run client.go
````
