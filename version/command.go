package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		Run:   printVersion,
	}
)

func GetVersion() string {
	return Version
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Println(GetVersion())
}
