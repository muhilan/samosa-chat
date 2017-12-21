package main

import (
	// "github.com/andlabs/ui"
	//"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/ui"
	"io/ioutil"
	"log"
	//"net/http"
	"os/user"
	"time"
	"bufio"
)

type ClientMetadata struct {
	Owner          string
	ChatServerHost string
	ChatServerPort string
	OwnerEmail     string
}


type MessageContext struct {
	Text string
	Owner string
	Time int64
}

var clientMD ClientMetadata
var home string

var messages = make(chan MessageContext)
var conn *tls.Conn

func main() {
	log.Print("Entry")

	clientMD = getClientMetaData()
    var err error
	conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%s", clientMD.ChatServerHost,clientMD.ChatServerPort), &tls.Config{InsecureSkipVerify : true})
	if err != nil {
		panic(err)
	}
	multi := ui.NewMultilineNonWrappingEntry()
	multi.ReadOnly()

	go func() {
		reader := bufio.NewReader(conn)
		for {
			in, err := reader.ReadString('\n')
			fmt.Printf("Received message %s\n",string(in))
			if err != nil {
				break
			}
			var msg MessageContext
			err = json.Unmarshal([]byte(in), &msg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			messages <- msg
		}
	}()
	go func(){
		select {
		case msg := <-messages:
			fmt.Println("About to paint message")
			multi.Append(formatText(msg.Time,msg.Text))
		}
	}()

	err = ui.Main(func() {
		newchat := ui.NewEntry()
		button := ui.NewButton("Send")

		box := ui.NewVerticalBox()
		horbox := ui.NewHorizontalBox()
		horbox.Append(newchat, true)
		horbox.Append(button, true)
		box.Append(multi, true)
		box.Append(horbox, false)

		window := ui.NewWindow("Samosa Chat", 300, 600, false)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			Post(newchat.Text(), time.Now().Unix())
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}


}

func formatText(epoch int64, newchat string) string {
	tm := time.Unix(epoch, 0)
	str := fmt.Sprintf("%s (%s): %s \n", tm.Format("2006-01-02 15:04:05"), clientMD.Owner, newchat)
	Post(str, epoch)
	return str
}

func Post(str string, epoch int64) {
	msgCtx := MessageContext{ Owner: clientMD.Owner, Time: epoch, Text: str }
	body := postPayload(createPayload(&msgCtx))
	fmt.Println(body)
	fmt.Fprint(conn, body)
}

/* Helper method to retrieve the client specific json file */

func getClientMetaData() ClientMetadata {
	raw, err := ioutil.ReadFile(homeDir() + "/.samosa-chat.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	var c ClientMetadata
	json.Unmarshal(raw, &c)
	return c
}

/* Util method to retrieve the home dir */
func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func formatDate(epoch int64) string {
    tm := time.Unix(epoch, 0)
    return tm.Format("2006-01-02 15:04:05")
}

func postPayload(msg string) string{
	return fmt.Sprintf("POST / HTTP/1.1\nContent-Type: application/json\nConnection: Keep-Alive\nHost: 127.0.0.1\nUser-Agent: chat-client\n\n%s\n\n",msg)
}

func createPayload(msgCtx *MessageContext) string {
	b, err := json.Marshal(msgCtx)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}