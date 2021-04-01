package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dundee/qrcode-payment-web/web"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var address string

var rootCmd = &cobra.Command{
	Use:   "qrcode-payment-web",
	Short: "Payment QR code generator",
	Long: `Payment QR code generator with web interface.

	Creates SPAYD and EPC QR codes`,
	Run: func(cmd *cobra.Command, args []string) {
		web.RunServer(address)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qrcode-payment-web.yaml)")
	rootCmd.PersistentFlags().StringVarP(&address, "bind-address", "b", "127.0.0.1:8080", "Address to bind to")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".qrcode-payment-web" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".qrcode-payment-web")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
