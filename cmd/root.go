package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/dlvhdr/turbo-compose/pkg/ui"
)

var (
	repository string

	rootCmd = &cobra.Command{
		Use: `turbo-compose --repository="something.amazonaws.com" path/to/docker-compose.yml`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return fmt.Errorf("please specify the docker-compose.yml path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			m := ui.NewModel(args[0], repository)
			p := tea.NewProgram(m)
			if _, err := p.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nil
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
