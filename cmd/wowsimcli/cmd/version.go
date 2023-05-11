package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "prints version information",
		Long:  "prints version information",
		Run: func(cmd *cobra.Command, args []string) {
			if version == "" {
				version = "development"
			}
			fmt.Println(version)
		},
	}
}
