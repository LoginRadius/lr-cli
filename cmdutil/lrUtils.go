package cmdutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

type LoginResponse struct {
	XSign   string `json:"xsign"`
	XToken  string `json:"xtoken"`
	AppName string `json:"app_name"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type APICred struct {
	Key    string `json:"Key"`
	Secret string `json:"Secret"`
}

func StoreCreds(cred *LoginResponse) error {
	user, _ := user.Current()

	os.Mkdir(filepath.Join(user.HomeDir, ".lrcli"), 0755)
	fileName := filepath.Join(user.HomeDir, ".lrcli", "token.json")

	dataBytes, _ := json.Marshal(cred)

	return ioutil.WriteFile(fileName, dataBytes, 0644)

}
func GetCreds() (*LoginResponse, error) {
	var v2 LoginResponse
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", "token.json")
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}

	file, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(file, &v2)
	return &v2, nil
}

func StoreAPICreds(cred *APICred) error {
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", "creds.json")
	dataBytes, _ := json.Marshal(cred)
	return ioutil.WriteFile(fileName, dataBytes, 0644)

}

func GetAPICreds() (*APICred, error) {
	var v APICred
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", "creds.json")
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}

	file, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(file, &v)
	return &v, nil
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
