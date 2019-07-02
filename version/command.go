package version

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		RunE: func(cmd *cobra.Command, args []string) error {
			verInfo := newVersionInfo()

			bz, err := json.MarshalIndent(verInfo, "", " ")
			if err != nil {
				return err
			}
			fmt.Println(string(bz))

			return nil
		},
	}

	return cmd
}
