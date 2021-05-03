package cmdutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
)

func ReadFile(filename string) ([]byte, error) {
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", filename)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}
	return ioutil.ReadFile(fileName)
}

func WriteFile(filename string, data []byte) error {
	user, _ := user.Current()
	os.Mkdir(filepath.Join(user.HomeDir, ".lrcli"), 0755)
	fileName := filepath.Join(user.HomeDir, ".lrcli", filename)
	return ioutil.WriteFile(fileName, data, 0644)
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

func DeleteFiles() error {
	user, _ := user.Current()
	dirName := filepath.Join(user.HomeDir, ".lrcli")
	dir, err := ioutil.ReadDir(dirName)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{dirName, d.Name()}...))
	}
	if err != nil {
		return err
	}
	return nil
}
