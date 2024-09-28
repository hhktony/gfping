package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var file string
var subnet string
var singleip string
var timeout int
var routinepool int
var output string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gfping",
	Short: "Batch network probe tool.",
	Long: `Batch network probe tool.`,
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
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gfping.yaml)")

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "read list of targets from a file")
	rootCmd.PersistentFlags().IntVarP(&routinepool, "concurrent", "c", 300, "number of goroutines to use (concurrent) (Default 300)")	
	rootCmd.PersistentFlags().StringVarP(&singleip, "singleip", "i", "", "single ip，EP: 192.168.1.1")
	rootCmd.PersistentFlags().StringVarP(&subnet, "subnet", "g", "", "generate target list (only if no -f -i specified), EP: 192.168.1.1/16")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 3000, "individual target initial timeout, unit ms")
	// rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "探测结果输出位置")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gfping" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gfping")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
