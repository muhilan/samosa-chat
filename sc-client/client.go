package main

import (
	"encoding/json"
	"fmt"
	"github.com/ProtonMail/ui"
	"io/ioutil"
	"log"
	"os/user"
	"time"
	"bufio"
	"net"
	"github.com/gen2brain/beeep"
	"os"
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

var messages = make(chan MessageContext)
var connectionC = make(chan net.Conn)
var conn net.Conn
var multi *ui.MultilineEntry

func main() {
	clientMD = getClientMetaData()
	go func() {
		select {
			case conn := <-connectionC :
			reader := bufio.NewReader(conn)
			for {
				in, err := reader.ReadString('\n')
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
		}
	}()
    var err error
	conn, err = net.Dial("tcp",  fmt.Sprintf("%s:%s", clientMD.ChatServerHost, clientMD.ChatServerPort))
	if err != nil {
		fmt.Println(err.Error())
	}
	conn.Write([]byte("\n"))
	connectionC <- conn

	go func(){
		for {
			select {
			case msg := <-messages:
				if msg.Text != "" {
					if msg.Owner != clientMD.Owner {
						beeep.Notify("New Message from " + msg.Owner, msg.Text)
						beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
					}
					if err != nil {
						panic(err)
					}
					ui.QueueMain(func () {
						multi.Append(formatText(msg.Time, msg.Text))
					})
				}
			}
		}
	}()

	err = ui.Main(func() {
		newChat := ui.NewEntry()
		button := ui.NewButton("Send")
		multi = ui.NewMultilineNonWrappingEntry()
		multi.ReadOnly()
		box := ui.NewVerticalBox()
		horbox := ui.NewHorizontalBox()
		horbox.Append(newChat, true)
		horbox.Append(button, true)
		box.Append(multi, true)
		box.Append(horbox, false)

		window := ui.NewWindow("Samosa Chat", 300, 600, true)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			post(newChat.Text(), time.Now().Unix())
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			os.Exit(0)
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}


}

func formatText(epoch int64, newChat string) string {
	tm := time.Unix(epoch, 0)
	str := fmt.Sprintf("%s (%s): %s \n", tm.Format("2006-01-02 15:04:05"), clientMD.Owner, newChat)
	return str
}

func post(str string, epoch int64) {
	msgCtx := MessageContext{ Owner: clientMD.Owner, Time: epoch, Text: str }
	body := createPayload(&msgCtx)
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

func createPayload(msgCtx *MessageContext) string {
	b, err := json.Marshal(msgCtx)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b) + "\n"
}