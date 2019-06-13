package cmd

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// Provides the commit identifier used to build the binary
var buildCode string

// Provides the UNIX timestamp of the build
var buildTimestamp string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		var components = map[string]string{
			"Version":    "0.1.0",
			"Build code": buildCode,
			"OS/Arch":    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			"Go version": runtime.Version(),
		}
		if buildTimestamp != "" {
			st, err := strconv.ParseInt(buildTimestamp, 10, 64)
			if err == nil {
				components["Release Date"] = time.Unix(st, 0).Format(time.RFC822)
			}
		}
		for k, v := range components {
			fmt.Printf("\033[21;37m%-13s:\033[0m %s\n", k, v)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
