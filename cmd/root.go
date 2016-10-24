package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dcfg",
	Short: "A tool for post-run configuration of Docker containers.",
	Long: `DrudConfig lets you store executeable info about how your app should be deployed
	alongside your app's source code.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "drud.yaml", "config file (default is ./drud.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
