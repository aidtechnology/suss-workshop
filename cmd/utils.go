package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/bryk-io/x/pki"
)

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
