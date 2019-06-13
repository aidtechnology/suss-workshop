package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/bryk-io/x/pki"
	"github.com/spf13/cobra"
)

var rootCACmd = &cobra.Command{
	Use:     "root-ca",
	Short:   "Create a root certificate authority",
	Example: "suss-workshop root-ca csr.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you need to provide the CSR json file")
		}
		csrJSON, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		cert, key, err := pki.RootCA(csrJSON)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile("root-ca.crt", cert, 0400); err != nil {
			return err
		}
		if err = ioutil.WriteFile("root-ca.pem", key, 0400); err != nil {
			return err
		}
		fmt.Println("root-ca created")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rootCACmd)
}
