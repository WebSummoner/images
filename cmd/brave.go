package cmd

import (
	"github.com/spf13/cobra"
	"github.com/websummoner/images/build"
)

var braveCmd = &cobra.Command{
	Use:   "brave",
	Short: "build Brave browser image",
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
		brave := &build.Brave{Requirements: req}
		return brave.Build()
	},
}
