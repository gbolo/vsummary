package command

import (
	"fmt"
	"os"

	"github.com/gbolo/vsummary/common"
	"github.com/gbolo/vsummary/poller"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	// pollendpointCmd represents the pollendpoint command
	pollendpointCmd = &cobra.Command{
		Use:   "pollendpoint",
		Short: "poll a specific vcenter endpoint via command line flags, then exit",
		Long: `Polls the specified vCenter then sends the data to the specified
vsummary-server URL for processing, then exits.

Example Usage:
  vsummary-poller pollendpoint -s vcenter.example.com -e HomeLab -username readonly 
`,
		Run: func(cmd *cobra.Command, args []string) {
			pollNow()
		},
	}

	// flag values
	flagVcenterHost, flagEnvironment, flagUsername, flagPassword string
)

func init() {
	rootCmd.AddCommand(pollendpointCmd)

	// flags
	pollendpointCmd.Flags().StringVarP(&flagVcenterHost, "vcenter", "s", "", "fqdn/ip of vcenter server")
	pollendpointCmd.Flags().StringVarP(&flagEnvironment, "environment", "e", "", "environment/name of vcenter (friendly name)")
	pollendpointCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "username for vcenter (readonly privilege needed)")
	pollendpointCmd.Flags().StringVarP(&flagPassword, "password", "p", "", "password for user (will prompt if not specified)")

	// mark required flags
	cobra.MarkFlagRequired(pollendpointCmd.Flags(), "vcenter")
	cobra.MarkFlagRequired(pollendpointCmd.Flags(), "environment")
	cobra.MarkFlagRequired(pollendpointCmd.Flags(), "username")
}

func pollNow() {

	// prompt for password if not set
	if flagPassword == "" {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Password for vCenter user %s", flagUsername),
			Mask:  '*',
			Validate: func(input string) error {
				if len(input) < 1 {
					return fmt.Errorf("password cannot be empty")
				}
				return nil
			},
		}
		var err error
		flagPassword, err = prompt.Run()
		if err != nil {
			fmt.Printf("Error with password: %v \n", err)
			os.Exit(1)
		}
	}

	// create external poller
	externalpoller := poller.NewExternalPoller(common.Poller{
		VcenterHost:       flagVcenterHost,
		VcenterName:       flagEnvironment,
		Username:          flagUsername,
		PlainTextPassword: flagPassword,
	})

	// set vsummary-server URL
	setVSummaryApiURL()
	err := externalpoller.SetApiUrl(vsummaryApiUrl)
	if err != nil {
		fmt.Printf("Error with vSummary Server URL: %v\n", err)
		os.Exit(1)
	}

	// poll then send results
	externalpoller.PollThenSend()
}
