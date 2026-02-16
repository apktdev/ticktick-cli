package cli

import (
	"context"

	"github.com/adrian/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func (a *app) newTasksCmd() *cobra.Command {
	tasks := &cobra.Command{Use: "tasks", Short: "Manage tasks"}

	var projectID string
	list := &cobra.Command{
		Use:   "list",
		Short: "List tasks in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(projectID, "--project-id"); err != nil {
				return err
			}
			data, err := a.client.ProjectData(context.Background(), projectID)
			if err != nil {
				return err
			}
			return a.print(data.Tasks)
		},
	}
	list.Flags().StringVar(&projectID, "project-id", "", "Project id")

	var title, content, due string
	var priority int
	var addProjectID string
	add := &cobra.Command{
		Use:   "add",
		Short: "Create a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(addProjectID, "--project-id"); err != nil {
				return err
			}
			if err := required(title, "--title"); err != nil {
				return err
			}
			dueDate, err := parseDueDate(due)
			if err != nil {
				return err
			}
			created, err := a.client.CreateTask(context.Background(), ticktick.Task{
				ProjectID: addProjectID,
				Title:     title,
				Content:   content,
				DueDate:   dueDate,
				Priority:  priority,
			})
			if err != nil {
				return err
			}
			return a.print(created)
		},
	}
	add.Flags().StringVar(&addProjectID, "project-id", "", "Project id")
	add.Flags().StringVar(&title, "title", "", "Task title")
	add.Flags().StringVar(&content, "content", "", "Task notes/content")
	add.Flags().StringVar(&due, "due", "", "Due date (RFC3339 or YYYY-MM-DD)")
	add.Flags().IntVar(&priority, "priority", 0, "Priority: 0 none, 1 low, 3 medium, 5 high")

	complete := &cobra.Command{
		Use:   "complete <project-id> <task-id>",
		Short: "Mark a task complete",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.client.CompleteTask(context.Background(), args[0], args[1]); err != nil {
				return err
			}
			return a.print("task completed")
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <project-id> <task-id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.client.DeleteTask(context.Background(), args[0], args[1]); err != nil {
				return err
			}
			return a.print("task deleted")
		},
	}

	tasks.AddCommand(list, add, complete, deleteCmd)
	return tasks
}
