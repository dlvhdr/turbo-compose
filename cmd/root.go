package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dlvhdr/turbo-compose/pkg/ui"
)

var (
	repository string

	rootCmd = &cobra.Command{
		Use: "turbo-compose",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := ui.NewModel(repository)
			return m.Run()
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "docker repository prefix")

}
