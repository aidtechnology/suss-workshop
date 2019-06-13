package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var userCertCmd = &cobra.Command{
	Use:     "user-cert",
	Short:   "Create a user certificate",
	Example: "suss-workshop user-cert user_csr.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read CSR
		if len(args) == 0 {
			return errors.New("you need to provide the CSR json file")
		}
		csrJSON, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		// Start CA instance
		ca, err := getCA()
		if err != nil {
			return err
		}

		// Sign user certificate
		cert, key, err := ca.SignRequestJSON(csrJSON, "user")
		if err != nil {
			return err
		}

		// Save client certificate
		if err = ioutil.WriteFile("user.crt", cert, 0400); err != nil {
			return err
		}
		if err = ioutil.WriteFile("user.pem", key, 0400); err != nil {
			return err
		}
		fmt.Println("user certificate created")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCertCmd)
}
