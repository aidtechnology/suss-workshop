package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCode = ""

var rootCmd = &cobra.Command{
	Use:   "suss-workshop",
	Short: "Sample application for the SUSS workshop of June 2019",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add support for ECHO_ env variables
	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("suss")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	})
}
