package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BuildVersion string = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(BuildVersion)
	},
}
