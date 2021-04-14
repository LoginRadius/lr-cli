package cmdutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

type APICred struct {
	Key    string `json:"Key"`
	Secret string `json:"Secret"`
}

func StoreCreds(cred []byte) error {
	user, _ := user.Current()

	os.Mkdir(filepath.Join(user.HomeDir, ".lrcli"), 0755)
	fileName := filepath.Join(user.HomeDir, ".lrcli", "token.json")

	return ioutil.WriteFile(fileName, cred, 0644)

}
func GetCreds() ([]byte, error) {
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", "token.json")
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}
	return ioutil.ReadFile(fileName)
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

func GeneratePassword() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var retVal string
	length := 0
	for length < 8 {
		retVal += string(charset[rand.Intn(len(charset))])
		length++
	}
	return retVal + "1@aA"

}

func ThemeConstants(theme string) (ThemeType, ThemeType) {
	auths := map[string]ThemeType{
		"London":   Theme1Auth,
		"Tokyo":    Theme2Auth,
		"Helsinki": Theme3Auth,
	}

	profiles := map[string]ThemeType{
		"London":   Theme1Profile,
		"Tokyo":    Theme2Profile,
		"Helsinki": Theme3Profile,
	}
	return auths[theme], profiles[theme]
}
