package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func (a *app) newProjectsCmd() *cobra.Command {
	projects := &cobra.Command{Use: "projects", Short: "Manage TickTick projects"}

	list := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := a.client.ListProjects(context.Background())
			if err != nil {
				return err
			}
			return a.print(items)
		},
	}

	projects.AddCommand(list)
	return projects
}
