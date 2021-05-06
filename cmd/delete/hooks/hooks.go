package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

type Delete struct {
	Isdeleted bool `json:"isdeleted"`
}

var hookid string
var option string

func NewHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "deletes hooks",
		Long: heredoc.Doc(`
		This command deletes webhooks configured with an App.
		`),
		Example: heredoc.Doc(`
			$ lr delete hooks --hookid <hookid>
			(Y)

			Webhook has been deleted.
 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if hookid == "" {
				return &cmdutil.FlagError{Err: errors.New("`--hookid` is required argument")}
			}
			return deleteHooks()
		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&hookid, "hookid", "i", "", "Hook Unique ID")
	return cmd
}

func deleteHooks() error {
	err := api.CurrentPlan()
	if err != nil {
		return err
	}
	checkHookID, err := api.CheckHookID(hookid)
	if err != nil {
		return err
	}
	if !checkHookID {
		fmt.Println("Hook ID does not exist.")
		return nil
	}
	fmt.Printf("Are you sure you want to proceed ? Press Y to continue: ")
	fmt.Scanf("%s", &option)
	if option != "Y" {
		return nil
	}
	isDeleted, err := delete()
	if err != nil {
		return err
	}
	if isDeleted {
		fmt.Println("Webhook has been deleted.")
	} else {
		fmt.Println("Delete action failed.")
	}
	return nil
}

func delete() (bool, error) {
	conf := config.GetInstance()
	delete := conf.AdminConsoleAPIDomain + "/integrations/webhook/" + hookid
	resp, err := request.Rest(http.MethodDelete, delete, nil, "")
	if err != nil {
		return false, err
	}
	var status Delete
	err = json.Unmarshal(resp, &status)
	if err != nil {
		return false, err
	}
	if status.Isdeleted == true {
		return true, nil
	}
	return false, nil

}
