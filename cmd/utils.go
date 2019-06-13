package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/bryk-io/x/did"
	"github.com/bryk-io/x/pki"
)

var tplUserCSR *template.Template

func init() {
	tplUserCSR, _ = template.New("csr").Parse(`{
  "cn": "{{.DID}}",
  "hosts": [
    "{{.DID}}"
  ],
  "key": {
    "algo": "ecdsa",
    "size": 521
  },
  "names": [
    {
      "o": "Singapore University of Social Sciences",
      "sa": "463 Clementi Road",
      "st": "Singapore",
      "pc": "599494",
      "c": "SG"
    }
  ]
}`)
}

func getCA() (*pki.CA, error) {
	confJSON, err := ioutil.ReadFile("ca_conf.json")
	if err != nil {
		return nil, errors.New("failed to read CA configuration file 'ca_conf.json'")
	}
	conf, err := pki.DecodeConfig(confJSON)
	if err != nil {
		return nil, err
	}
	return pki.NewCA("root-ca.crt", "root-ca.pem", nil, conf)
}

func resolveDID(value string) (*did.Identifier, error) {
	// Verify the provided value is a valid DID string
	d, err := did.Parse(value)
	if err != nil {
		return nil, err
	}
	if d.Method() != "bryk" {
		return nil, errors.New("only the 'bryk' DID method is supported for this application")
	}

	// Retrieve element
	res, err := http.Get(fmt.Sprintf("https://did.bryk.io/v1/retrieve?subject=%s", d.Subject()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse document
	docJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	doc := &did.Document{}
	if err = json.Unmarshal(docJSON, doc); err != nil {
		return nil, err
	}
	return did.FromDocument(doc)
}

func verifySignature(id *did.Identifier, challenge string, sig *did.SignatureLD) error {
	masterKey := id.Key("master")
	if masterKey == nil {
		return errors.New("failed to retrieve master key for the DID")
	}
	if !masterKey.VerifySignatureLD([]byte(challenge), sig) {
		return errors.New("invalid signature/challenge")
	}
	return nil
}
