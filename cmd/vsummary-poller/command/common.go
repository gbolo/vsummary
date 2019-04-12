package command

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"os"
)

var vsummaryApiUrl string

func setVSummaryApiURL() {
	// check viper/cobra
	vsummaryApiUrl = viper.GetString("poller.url")

	// do a prompt if not set
	if vsummaryApiUrl == "" {
		fmt.Printf("vSummary Server URL is not set in config or flag! please enter this below\n")
		prompt := promptui.Prompt{
			Label: "vSummary server base URL",
			Validate: func(input string) error {
				if len(input) < 1 {
					return fmt.Errorf("URL cannot be empty")
				}
				return nil
			},
		}
		var err error
		vsummaryApiUrl, err = prompt.Run()
		if err != nil {
			fmt.Printf("Error with url: %v \n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Setting vSummary Server Base URL: %s\n", vsummaryApiUrl)
}
