package cmd

import (
	"github.com/spf13/cobra"
	"github.com/websummoner/images/build"
)

var (
	firefoxCmd = &cobra.Command{
		Use:   "firefox",
		Short: "build Firefox image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := build.Requirements{
				BrowserSource:  build.BrowserSource(browserSource),
				BrowserChannel: browserChannel,
				DriverVersion:  driverVersion,
				NoCache:        noCache,
				TestsDir:       testsDir,
				RunTests:       test,
				IgnoreTests:    ignoreTests,
				Tags:           tags,
				PushImage:      push,
			}
			firefox := &build.Firefox{Requirements: req}
			return firefox.Build()
		},
	}
)

func init() {
}
