package cmd

import (
	"fmt"
	"os"
	"quebic-faas/quebic-faas-cli/app"

	"github.com/spf13/cobra"
)

var appContainer app.App

func init() {
	setupCmds()
	setupFlags()
}

func setupCmds() {
	rootCmd.AddCommand(functionCmd)
	rootCmd.AddCommand(routeCmd)
	rootCmd.AddCommand(requestTrackerCmd)
	rootCmd.AddCommand(mgrCompCmd)
	rootCmd.AddCommand(configCmd)
}

func setupFlags() {

}

var rootCmd = &cobra.Command{
	Use:   "quebic",
	Short: "quebic is a development kit for write serverless function",
	Long:  `quebic is a development kit for write serverless function. Complete documentation is available at http://quebic.io`,
}

//Execute execute cmd
func Execute() {

	appContainer = app.App{}
	appContainer.Start()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
