package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bryk-io/x/pki"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a server instance for the sample digital service",
	RunE:  runServer,
}

type serviceResponse struct {
	Ok bool `json:"ok"`
	Response interface{} `json:"response"`
}

func (sr *serviceResponse) encode() []byte {
	r, _ := json.Marshal(sr)
	return r
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, _ []string) error {
	// Get server's certificate authority
	ca, err := getCA()
	if err != nil {
		return err
	}

	// Setup server's router
	router := mux.NewRouter()
	router.HandleFunc("/enroll", enrollHandler(ca)).Methods(http.MethodPost)
	router.HandleFunc("/connect", connectHandler(ca)).Methods(http.MethodGet)
	router.PathPrefix("/").HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		r := &serviceResponse{
			Ok: true,
			Response: "SUSS workshop sample service =D",
		}
		res.Write(r.encode())
	})

	// Start server
	srv := &http.Server{
		Handler:      router,
		Addr:         ":9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("server ready")
	fmt.Println("waiting for connections at port: 9090")
	return srv.ListenAndServe()
}

// Enroll
// The server will generate a client certificate for the user.
// Enrollment requests include the following fields:
// - did: subject's DID to use
// - challenge: a random value generated to authorize the transaction
// - signature: signature generated for the challenge
//
// To process the enrollment the server performs the following:
// - Resolve the DID
// - Verify the signature/challenge is valid
// - Generate a certificate and private key for the DID
func enrollHandler(ca *pki.CA) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		r := &serviceResponse{
			Ok: false,
			Response: "",
		}

		// Read request body
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil || len(body) == 0 {
			res.WriteHeader(400)
			r.Response = "empty request"
			res.Write(r.encode())
			return
		}

		// Decode enrollment request
		er := &enrollmentRequest{}
		if err = json.Unmarshal(body, er); err != nil {
			res.WriteHeader(400)
			r.Response = "invalid request contents"
			res.Write(r.encode())
			return
		}

		// Resolve provided DID
		id, err := resolveDID(er.Did)
		if err != nil {
			res.WriteHeader(400)
			r.Response = "failed to resolve DID"
			res.Write(r.encode())
			return
		}

		// Validate challenge/signature
		if err = verifySignature(id, er.Challenge, er.Signature); err != nil {
			res.WriteHeader(400)
			r.Response = err.Error()
			res.Write(r.encode())
			return
		}

		// Generate CSR
		buf := bytes.NewBuffer(nil)
		if err = tplUserCSR.Execute(buf, map[string]string{"DID":id.String()}); err != nil {
			res.WriteHeader(400)
			r.Response = "failed to generate CSR"
			res.Write(r.encode())
			return
		}

		// Generate certificate
		cert, key, err := ca.SignRequestJSON(buf.Bytes(), "user")
		if err != nil {
			res.WriteHeader(400)
			r.Response = "failed to generate certificate"
			res.Write(r.encode())
			return
		}

		// All good!
		r.Ok = true
		r.Response = &enrollmentResponse{
			Cert: cert,
			Key:  key,
		}
		res.Write(r.encode())
		return
	}
}

// Connect
// Receive a user request to start a session with the service.
// The server will validate the client certificate to prevent unauthorized access.
func connectHandler(ca *pki.CA) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		r := &serviceResponse{
			Ok: true,
			Response: "connecting",
		}
		res.Write(r.encode())
	}
}
