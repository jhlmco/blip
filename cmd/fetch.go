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

// pullCmd represents the pull command
var fetchCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if path == "" {
			err := fmt.Errorf("failed to provid path")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}
		fetch.Fetch(path, filename)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&path, "path", "p", "", "Path to artifact.  Can be URL or Local Path.")
	fetchCmd.Flags().StringVarP(&filename, "file", "f", "", "Filename for artifact.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
