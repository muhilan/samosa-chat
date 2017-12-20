package main

import (
    // "github.com/andlabs/ui"
    "log"
    "github.com/ProtonMail/ui"
    "time"
    "net/http"
    "fmt"
    "crypto/tls"
    "io/ioutil"
    "bytes"
    "encoding/json"
    "os/user"
)
type ClientMetadata struct{
    Owner string
    ChatServerHost string
    ChatServerPort string
    OwnerEmail string
}
type Message struct {
    Owner   string  
    Time string 
    Text string 
}

type MessageContext struct {
    msgs []Message  `json:"msgs"`
}

var clientMD ClientMetadata
var home string

func main() {
    log.Print("Entry")
    clientMD = getClientMetaData()
    err := ui.Main(func() {
        newchat := ui.NewEntry()
        button := ui.NewButton("Send")
        multi := ui.NewMultilineNonWrappingEntry()
        multi.ReadOnly()
        box := ui.NewVerticalBox()
        horbox := ui.NewHorizontalBox()
        horbox.Append(newchat, true)
        horbox.Append(button, true)
        box.Append(multi,true)
        box.Append(horbox, false)
       
        window := ui.NewWindow("Samosa Chat", 300, 600, false)
        window.SetChild(box)
        button.OnClicked(func (*ui.Button) {
            multi.Append(formatText(newchat.Text()))
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

func formatText(newchat string) string{
    currentEpocy := int64(time.Now().Unix())
    tm := time.Unix(currentEpocy, 0)
    fmt.Println(tm)

    str := fmt.Sprintf(  "%s (%s): %s \n", tm.Format("2006-01-02 15:04:05"), clientMD.Owner, newchat)
    Post(str)
    return  str
}

func Post(str string){

    var msg = &Message { Owner: "mgm", Time : time.Now().Format("2006-01-02 15:04:05"), Text: str}
    jsonStr, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err)
        return
    }

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    r, _ := http.NewRequest("POST", "https://localhost:8080", bytes.NewBuffer(jsonStr)) // <-- URL-encoded payload
    r.Header.Add("Content-Type", "application/json")

    resp, err := client.Do(r)
    if err!=nil {
        log.Print(err)
    } else {
        fmt.Println(resp.Status)
        data, _ := ioutil.ReadAll(resp.Body)
        fmt.Println(string(data))
        resp.Body.Close()
    }

}

/* Helper method to retrieve the client specific json file */

func getClientMetaData() ClientMetadata {
    raw, err := ioutil.ReadFile(homeDir()+ "/.samosa-chat.json")
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
        log.Fatal( err )
    }
    return usr.HomeDir
}