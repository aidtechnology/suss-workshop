package cmd

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bryk-io/x/cli"
	"github.com/chzyer/readline"
	"github.com/gorilla/websocket"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the sample digital service",
	RunE:  runClient,
}

type session struct {
	alias   string
	ws      *websocket.Conn
	rl      *readline.Instance
	errChan chan error
}

func (s *session) readConsole() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			s.errChan <- err
			return
		}

		if line == "bye" || line == "close" || line == "exit" {
			msg := fmt.Sprintf("%s: %s", aurora.Red(s.alias), "is going away")
			s.ws.WriteMessage(websocket.TextMessage, []byte(msg))
			s.errChan <- nil
			return
		}

		line = fmt.Sprintf("%s: %s", s.alias, line)
		if err = s.ws.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			s.errChan <- err
			return
		}
	}
}

func (s *session) readWebsocket() {
	for {
		msgType, buf, err := s.ws.ReadMessage()
		if err != nil {
			s.errChan <- err
			return
		}
		var text string
		switch msgType {
		case websocket.TextMessage:
			text = string(buf)
		default:
			s.errChan <- fmt.Errorf("unknown websocket frame type: %d", msgType)
			return
		}
		segs := strings.Split(text, ":")
		if segs[0] == s.alias {
			fmt.Fprint(s.rl.Stdout(), fmt.Sprintf("%s: %s\n", aurora.Yellow(segs[0]), segs[1]))
		} else {
			fmt.Fprint(s.rl.Stdout(), fmt.Sprintf("%s: %s\n", aurora.Blue(segs[0]), segs[1]))
		}
	}
}

func init() {
	name, err := os.Hostname()
	if err != nil {
		name = fmt.Sprintf("user-%d", time.Now().Second())
	}
	params := []cli.Param{
		{
			Name:      "cert",
			Usage:     "certificate to access the service",
			FlagKey:   "connect.cert",
			ByDefault: "",
		},
		{
			Name:      "alias",
			Usage:     "alias for the session",
			FlagKey:   "connect.alias",
			ByDefault: name,
		},
	}
	if err := cli.SetupCommandParams(connectCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(connectCmd)
}

func runClient(cmd *cobra.Command, args []string) error {
	// Validate parameters
	if len(args) == 0 {
		return errors.New("you need to specify the service endpoint")
	}
	if viper.GetString("connect.cert") == "" {
		return errors.New("you need to provide your user certificate")
	}

	// Load certificate
	c, err := ioutil.ReadFile(viper.GetString("connect.cert"))
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("ws://%s/connect", args[0])
 	headers := make(http.Header)
 	headers.Set("X-user-certificate", base64.StdEncoding.EncodeToString(c))
	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig:&tls.Config{
			InsecureSkipVerify: true,
		},
	}
	ws, _, err := dialer.Dial(endpoint, headers)
	if err != nil {
		return err
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("%s", aurora.Magenta("» ")),
		InterruptPrompt: "^C",
		EOFPrompt:       "bye",
	})
	if err != nil {
		return err
	}
	defer rl.Close()
	sess := &session{
		alias:   viper.GetString("connect.alias"),
		ws:      ws,
		rl:      rl,
		errChan: make(chan error),
	}
	go sess.readConsole()
	go sess.readWebsocket()
	return <-sess.errChan
}
