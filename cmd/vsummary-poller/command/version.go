package command

import (
	"github.com/gbolo/vsummary/common"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Outputs version information",
	Long: `Outputs the following version information:
  version     : version of code
  build date  : date that this binary was built on
  git hash    : short commit hash from git repo
  go version  : version of go used to compile this binary
  go compiler : type of go compiler used
  platform    : OS/Arch of platform this binary is compiled for
`,
	Run: func(cmd *cobra.Command, args []string) {
		common.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
