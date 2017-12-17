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
)

type Message struct {
    Owner   string  
    Time string 
    Text string 
}

type MessageContext struct {
    msgs []Message  `json:"msgs"`
}

func main() {
    log.Print("Entry")
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
    str := string(time.Now().Format("2006-01-02 15:04:05")) + "  " + newchat + "\n"
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

