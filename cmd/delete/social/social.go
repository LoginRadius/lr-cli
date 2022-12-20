package social

import (
	"errors"
	"fmt"
	"strings"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/api"

	"github.com/spf13/cobra"
)

type provider struct {
	ProviderName string `json:"providerName"`
}

type Result struct {
	Isdeleted bool `json:"isdeleted"`
}

var url string

func NewsocialCmd() *cobra.Command {
	opts := &provider{}

	cmd := &cobra.Command{
		Use:   "social",
		Short: "Deletes a social provider",
		Long:  `Use this command to delete a configured social login provider from your application.`,
		Example: `
		$ lr delete social -p Google
		? Are you Sure you want to delete the provider? Yes
		Google Successfully Deleted
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.ProviderName == "" {
				return &cmdutil.FlagError{Err: errors.New("`provider` is require argument")}
			}
			return delete(opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.ProviderName, "provider", "p", "", "Enter name of the provider you want to delete")

	return cmd
}

func delete(opts *provider) error {

	var shouldDelete bool
	if err := prompt.Confirm("Are you Sure you want to delete the provider?", &shouldDelete); err != nil {
		return err
	}

	if shouldDelete {
		var validProvider bool
		activeProv, err := api.GetActiveProviders()
		for  prov, _ := range activeProv {
			ok := prov == strings.ToLower(opts.ProviderName)
			if ok {
				validProvider = true
			}
		}
		if validProvider {
		 err = api.UpdateProviderStatus(opts.ProviderName, false)
			if err != nil {
				return err
			}
		
		
			fmt.Println(opts.ProviderName + " Successfully Deleted")
		} else {
			return &cmdutil.FlagError{Err: errors.New(opts.ProviderName + " not added as a social provider")}
		} 
		
	}
	return nil
}
