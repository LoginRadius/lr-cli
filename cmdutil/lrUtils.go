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

	"github.com/shirou/gopsutil/host"

	"github.com/loginradius/lr-cli/internal/build"
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
		"Template_1":   Theme1Auth,
		"Template_2":   Theme2Auth,
		"Template_3": 	Theme3Auth,
		"Template_4": 	Theme4Auth,
		"Template_5": 	Theme5Auth,
	}

	profiles := map[string]ThemeType{
		"Template_1":   Theme1Profile,
		"Template_2":   Theme2Profile,
		"Template_3": 	Theme3Profile,
		"Template_4": 	Theme4Profile,
		"Template_5": 	Theme5Profile,
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

func DeleteFile(filename string) error {
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", filename)
	err := os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}
func UAString() string {
	if info, err := host.Info(); err == nil {
		return `LRCLI/` + build.Version + " (" + info.Platform + "; " + info.KernelArch + ") " + info.OS + "/" + info.PlatformVersion + " (" + info.Hostname + ")"
	}
	return ""
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
