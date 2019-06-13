package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bryk-io/x/cli"
	"github.com/bryk-io/x/did"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type enrollmentRequest struct {
	Did       string           `json:"did"`
	Challenge string           `json:"challenge"`
	Signature *did.SignatureLD `json:"signature"`
}

type enrollmentResponse struct {
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Enroll a given DID with the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate parameters
		log.Println("validating parameters...")
		userDID := viper.GetString("enroll.did")
		if userDID == "" {
			return errors.New("DID value is required")
		}
		challenge := viper.GetString("enroll.challenge")
		if challenge == "" {
			return errors.New("a challenge value is required")
		}
		signature := viper.GetString("enroll.signature")
		if signature == "" {
			return errors.New("a signature file is required")
		}
		endpoint := viper.GetString("enroll.endpoint")
		if endpoint == "" {
			return errors.New("you need to specify the service endpoint to use")
		}

		// Resolve DID
		log.Println("retrieving DID...")
		id, err := resolveDID(userDID)
		if err != nil {
			return err
		}

		// Load signature
		log.Println("loading challenge signature...")
		sigJSON, err := ioutil.ReadFile(signature)
		if err != nil {
			return err
		}
		sigLD := &did.SignatureLD{}
		if err = json.Unmarshal(sigJSON, sigLD); err != nil {
			return err
		}

		// Verify challenge on the client side
		log.Println("verifying challenge signature...")
		if err = verifySignature(id, challenge, sigLD); err != nil {
			return err
		}

		// Submit enrollment request
		log.Println("submitting enrollment request...")
		req := &enrollmentRequest{
			Did:       id.String(),
			Challenge: challenge,
			Signature: sigLD,
		}
		js, _ := json.MarshalIndent(req, "", "  ")
		fmt.Printf("%s", js)
		res, err := http.Post(fmt.Sprintf("%s/enroll", endpoint), "application/json", bytes.NewReader(js))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// Parse service response
		log.Println("inspecting service response...")
		sr := &serviceResponse{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(body, sr); err != nil {
			return err
		}
		if !sr.Ok {
			return errors.New(sr.Response.(string))
		}
		log.Println("saving obtained certificate...")
		creds := sr.Response.(map[string]interface{})
		cert, _ := base64.StdEncoding.DecodeString(creds["cert"].(string))
		key, _ := base64.StdEncoding.DecodeString(creds["key"].(string))
		if err = ioutil.WriteFile(fmt.Sprintf("%s.crt", id.Subject()), cert, 0400); err != nil {
			return err
		}
		if err = ioutil.WriteFile(fmt.Sprintf("%s.pem", id.Subject()), key, 0400); err != nil {
			return err
		}
		log.Println("certificate saved successfully!")
		return nil
	},
}

func init() {
	params := []cli.Param{
		{
			Name:      "did",
			Usage:     "DID to enroll with the service",
			FlagKey:   "enroll.did",
			ByDefault: "",
		},
		{
			Name:      "challenge",
			Usage:     "challenge value used for authentication during the enrollment process",
			FlagKey:   "enroll.challenge",
			ByDefault: "",
		},
		{
			Name:      "signature",
			Usage:     "signature produced for authentication",
			FlagKey:   "enroll.signature",
			ByDefault: "",
		},
		{
			Name:      "endpoint",
			Usage:     "service endpoint to send the enrollment request to",
			FlagKey:   "enroll.endpoint",
			ByDefault: "",
		},
	}
	if err := cli.SetupCommandParams(enrollCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(enrollCmd)
}
