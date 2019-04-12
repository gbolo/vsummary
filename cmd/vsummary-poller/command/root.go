package command

import (
	"fmt"
	"os"

	"github.com/gbolo/vsummary/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vsummary-poller",
	Short: "A poller for vcenter server(s) which sends data back to specified vsummary-server",
	Long: `This poller can collect data from one or more vcenter server(s) and send that data back
to a vsummary-server URL for processing. This poller offers two modes:
 - pollendpoint: polls a specific endpoint once then exits.
 - daemonize:    polls a list of vcenter servers defined in the configuration file at a desired interval.

**NOTE** the credential for the vCenter user requires READ-ONLY access to the top level vCenter object`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags for our cli
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().String("log-level", "INFO", "supported levels: INFO, WARNING, CRITICAL, DEBUG")
	rootCmd.PersistentFlags().String("vsummary-url", "", "vsummary-server URL")

	// viper integration
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("poller.url", rootCmd.PersistentFlags().Lookup("vsummary-url"))
}

// initConfig calls the usual vSummary config init
func initConfig() {
	config.ConfigInitPoller(cfgFile)
}
