package logout

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/spf13/cobra"
)

var tempServer *cmdutil.TempServer

func NewLogoutCmd() *cobra.Command {
	conf := config.GetInstance()
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout of LR account",
		Long:  `This commmand logs user out of the LR account`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdutil.Openbrowser(conf.HubPageDomain + "/auth.aspx?action=logout&return_url=http://localhost:8089/postLogout")
			tempServer = cmdutil.CreateTempServer(cmdutil.TempServer{
				Port:        ":8089",
				HandlerFunc: postLogout,
				RouteName:   "/postLogout",
			})
			tempServer.Server.ListenAndServe()
			return logout()
		},
	}
	return cmd
}

func logout() error {
	user, _ := user.Current()
	fileName := filepath.Join(user.HomeDir, ".lrcli", "token.json")
	dirName := filepath.Join(user.HomeDir, ".lrcli")
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return &cmdutil.FlagError{Err: errors.New(" You have already been logged Out")}
	} else {
		dir, err := ioutil.ReadDir(dirName)
		for _, d := range dir {
			os.RemoveAll(path.Join([]string{dirName, d.Name()}...))
		}
		if err != nil {
			return err
		}

	}
	log.Println("You are successfully Logged Out")
	return nil
}

func postLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are Successfully Logged Out, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}
