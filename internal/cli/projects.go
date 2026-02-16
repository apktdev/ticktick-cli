package cli

import (
	"context"

	"github.com/adrian/ticktick-cli/internal/ticktick"
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

	get := &cobra.Command{
		Use:   "get <project-id>",
		Short: "Get project by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := a.client.GetProject(context.Background(), args[0])
			if err != nil {
				return err
			}
			return a.print(item)
		},
	}

	var name, color, viewMode, kind string
	var sortOrder int64
	add := &cobra.Command{
		Use:   "add",
		Short: "Create project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(name, "--name"); err != nil {
				return err
			}
			item, err := a.client.CreateProject(context.Background(), ticktick.Project{
				Name:      name,
				Color:     color,
				SortOrder: sortOrder,
				ViewMode:  viewMode,
				Kind:      kind,
			})
			if err != nil {
				return err
			}
			return a.print(item)
		},
	}
	add.Flags().StringVar(&name, "name", "", "Project name")
	add.Flags().StringVar(&color, "color", "", `Project color (e.g. "#F18181")`)
	add.Flags().Int64Var(&sortOrder, "sort-order", 0, "Project sort order")
	add.Flags().StringVar(&viewMode, "view-mode", "", `View mode: "list", "kanban", "timeline"`)
	add.Flags().StringVar(&kind, "kind", "", `Project kind: "TASK" or "NOTE"`)

	var updateName, updateColor, updateViewMode, updateKind string
	var updateSortOrder int64
	update := &cobra.Command{
		Use:   "update <project-id>",
		Short: "Update project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := a.client.UpdateProject(context.Background(), args[0], ticktick.Project{
				Name:      updateName,
				Color:     updateColor,
				SortOrder: updateSortOrder,
				ViewMode:  updateViewMode,
				Kind:      updateKind,
			})
			if err != nil {
				return err
			}
			return a.print(item)
		},
	}
	update.Flags().StringVar(&updateName, "name", "", "Project name")
	update.Flags().StringVar(&updateColor, "color", "", `Project color (e.g. "#F18181")`)
	update.Flags().Int64Var(&updateSortOrder, "sort-order", 0, "Project sort order")
	update.Flags().StringVar(&updateViewMode, "view-mode", "", `View mode: "list", "kanban", "timeline"`)
	update.Flags().StringVar(&updateKind, "kind", "", `Project kind: "TASK" or "NOTE"`)

	deleteCmd := &cobra.Command{
		Use:   "delete <project-id>",
		Short: "Delete project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.client.DeleteProject(context.Background(), args[0]); err != nil {
				return err
			}
			return a.print("project deleted")
		},
	}

	projects.AddCommand(list, get, add, update, deleteCmd)
	return projects
}
