/*
Copyright © 2023 Jerod Heck <jerod.heck@lmco.com>
Software Factory InnerSource License
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jhlmco/blip/pkg/fetch"
	"github.com/spf13/cobra"
)

var path = ""
var filename = ""
var username string
var password string

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Command to fetch file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if path == "" {
			err := fmt.Errorf("failed to provid path or sbom")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}
		fetch.Fetch(path, filename, username, password)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&path, "path", "", "", "Path to artifact.  Can be URL")
	fetchCmd.Flags().StringVarP(&filename, "file", "f", "", "Filename for artifact.")
	fetchCmd.Flags().StringVarP(&username, "user", "u", "", "Username for authentication.")
	fetchCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication.")

}
