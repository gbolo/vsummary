package command

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gbolo/vsummary/poller"
	"github.com/spf13/cobra"
)

// daemonizeCmd represents the daemonize command
var daemonizeCmd = &cobra.Command{
	Use:   "daemonize",
	Short: "Daemonize this poller for use as a service",
	Long: `When daemonizing a poller, an interval is set (in minutes) and this poller
will remain running as a proccess polling all vcenter(s) defined in the supplied
configuration file indefinitely`,
	Run: func(cmd *cobra.Command, args []string) {
		daemonize()
	},
}

func init() {
	rootCmd.AddCommand(daemonizeCmd)
}

func daemonize() {

	// set vsummary-server URL
	setVSummaryApiURL()

	// get pollers defined in configuration file
	externalPollers := poller.GetExternalPollersFromConfig()
	if len(externalPollers) < 1 {
		fmt.Printf("Error: did not find any pollers defined in configuration file!")
		os.Exit(1)
	}

	// daemonize all pollers
	for i := range externalPollers {
		err := externalPollers[i].SetApiUrl(vsummaryApiUrl)
		if err != nil {
			fmt.Printf("Error with vSummary Server URL: %v\n", err)
			os.Exit(1)
		}
		go externalPollers[i].Daemonize()

		// sleep to stagger polling a bit
		time.Sleep(5 * time.Second)
	}

	// block forever
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
