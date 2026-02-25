package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/apktdev/ticktick-cli/internal/ticktick"
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

	get := &cobra.Command{
		Use:   "get <project-id> <task-id>",
		Short: "Get task by project and task id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := a.client.GetTask(context.Background(), args[0], args[1])
			if err != nil {
				return err
			}
			return a.print(item)
		},
	}

	var addProjectID, addTitle, addContent, addDesc, addStart, addDue, addTimeZone, addRepeat, addKind string
	var addPriority int
	var addAllDay bool
	var addReminders []string
	add := &cobra.Command{
		Use:   "add",
		Short: "Create a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(addProjectID, "--project-id"); err != nil {
				return err
			}
			if err := required(addTitle, "--title"); err != nil {
				return err
			}
			startDate, err := parseDateTime(addStart)
			if err != nil {
				return err
			}
			dueDate, err := parseDateTime(addDue)
			if err != nil {
				return err
			}
			kind, err := normalizeTaskKind(addKind)
			if err != nil {
				return err
			}
			created, err := a.client.CreateTask(context.Background(), ticktick.Task{
				ProjectID:  addProjectID,
				Title:      addTitle,
				Content:    addContent,
				Desc:       addDesc,
				StartDate:  startDate,
				DueDate:    dueDate,
				TimeZone:   addTimeZone,
				RepeatFlag: addRepeat,
				Reminders:  addReminders,
				Priority:   addPriority,
				IsAllDay:   addAllDay,
				Kind:       kind,
			})
			if err != nil {
				return err
			}
			return a.print(created)
		},
	}
	add.Flags().StringVar(&addProjectID, "project-id", "", "Project id")
	add.Flags().StringVar(&addTitle, "title", "", "Task title")
	add.Flags().StringVar(&addContent, "content", "", "Task content")
	add.Flags().StringVar(&addDesc, "desc", "", "Task checklist description")
	add.Flags().StringVar(&addStart, "start", "", "Start date (RFC3339 or YYYY-MM-DD)")
	add.Flags().StringVar(&addDue, "due", "", "Due date (RFC3339 or YYYY-MM-DD)")
	add.Flags().StringVar(&addTimeZone, "time-zone", "", "IANA timezone (e.g. America/Los_Angeles)")
	add.Flags().StringVar(&addRepeat, "repeat", "", "Recurring rule (e.g. RRULE:FREQ=DAILY;INTERVAL=1)")
	add.Flags().StringVar(&addKind, "kind", "", "Item kind: TEXT, NOTE, CHECKLIST")
	add.Flags().StringSliceVar(&addReminders, "reminder", nil, "Reminder trigger(s), repeatable")
	add.Flags().BoolVar(&addAllDay, "all-day", false, "Mark task as all-day")
	add.Flags().IntVar(&addPriority, "priority", 0, "Priority: 0 none, 1 low, 3 medium, 5 high")

	var updateProjectID, updateTitle, updateContent, updateDesc, updateStart, updateDue, updateTimeZone, updateRepeat, updateKind string
	var updatePriority int
	var updateAllDay bool
	var updateReminders []string
	update := &cobra.Command{
		Use:   "update <task-id>",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(updateProjectID, "--project-id"); err != nil {
				return err
			}
			startDate, err := parseDateTime(updateStart)
			if err != nil {
				return err
			}
			dueDate, err := parseDateTime(updateDue)
			if err != nil {
				return err
			}
			kind, err := normalizeTaskKind(updateKind)
			if err != nil {
				return err
			}
			updated, err := a.client.UpdateTask(context.Background(), args[0], ticktick.Task{
				ID:         args[0],
				ProjectID:  updateProjectID,
				Title:      updateTitle,
				Content:    updateContent,
				Desc:       updateDesc,
				StartDate:  startDate,
				DueDate:    dueDate,
				TimeZone:   updateTimeZone,
				RepeatFlag: updateRepeat,
				Reminders:  updateReminders,
				Priority:   updatePriority,
				IsAllDay:   updateAllDay,
				Kind:       kind,
			})
			if err != nil {
				return err
			}
			return a.print(updated)
		},
	}
	update.Flags().StringVar(&updateProjectID, "project-id", "", "Project id (required by API)")
	update.Flags().StringVar(&updateTitle, "title", "", "Task title")
	update.Flags().StringVar(&updateContent, "content", "", "Task content")
	update.Flags().StringVar(&updateDesc, "desc", "", "Task checklist description")
	update.Flags().StringVar(&updateStart, "start", "", "Start date (RFC3339 or YYYY-MM-DD)")
	update.Flags().StringVar(&updateDue, "due", "", "Due date (RFC3339 or YYYY-MM-DD)")
	update.Flags().StringVar(&updateTimeZone, "time-zone", "", "IANA timezone (e.g. America/Los_Angeles)")
	update.Flags().StringVar(&updateRepeat, "repeat", "", "Recurring rule (e.g. RRULE:FREQ=DAILY;INTERVAL=1)")
	update.Flags().StringVar(&updateKind, "kind", "", "Item kind: TEXT, NOTE, CHECKLIST")
	update.Flags().StringSliceVar(&updateReminders, "reminder", nil, "Reminder trigger(s), repeatable")
	update.Flags().BoolVar(&updateAllDay, "all-day", false, "Mark task as all-day")
	update.Flags().IntVar(&updatePriority, "priority", 0, "Priority: 0 none, 1 low, 3 medium, 5 high")

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

	tasks.AddCommand(list, get, add, update, complete, deleteCmd)
	return tasks
}

func normalizeTaskKind(kind string) (string, error) {
	if kind == "" {
		return "", nil
	}
	v := strings.ToUpper(strings.TrimSpace(kind))
	switch v {
	case "TEXT", "NOTE", "CHECKLIST":
		return v, nil
	default:
		return "", fmt.Errorf("invalid --kind %q; use TEXT, NOTE, or CHECKLIST", kind)
	}
}
