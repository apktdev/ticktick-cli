package cli

import (
	"context"
	"encoding/json"
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

	var addProjectID, addTitle, addContent, addDesc, addStart, addDue, addTimeZone, addRepeat, addKind, addItemsJSON string
	var addPriority int
	var addSortOrder int64
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
			items, err := parseChecklistItems(addItemsJSON)
			if err != nil {
				return err
			}

			payload := ticktick.Task{
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
				SortOrder:  addSortOrder,
				IsAllDay:   addAllDay,
				Kind:       kind,
			}
			if items != nil {
				payload.Items = items
			}

			created, err := a.client.CreateTask(context.Background(), payload)
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
	add.Flags().StringVar(&addItemsJSON, "items-json", "", "Checklist items JSON array")
	add.Flags().BoolVar(&addAllDay, "all-day", false, "Mark task as all-day")
	add.Flags().IntVar(&addPriority, "priority", 0, "Priority: 0 none, 1 low, 3 medium, 5 high")
	add.Flags().Int64Var(&addSortOrder, "sort-order", 0, "Task sort order")

	var updateProjectID, updateTitle, updateContent, updateDesc, updateStart, updateDue, updateTimeZone, updateRepeat, updateKind, updateItemsJSON string
	var updatePriority int
	var updateSortOrder int64
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
			if !cmd.Flags().Changed("title") && !cmd.Flags().Changed("content") && !cmd.Flags().Changed("desc") &&
				!cmd.Flags().Changed("start") && !cmd.Flags().Changed("due") && !cmd.Flags().Changed("time-zone") &&
				!cmd.Flags().Changed("repeat") && !cmd.Flags().Changed("kind") && !cmd.Flags().Changed("reminder") &&
				!cmd.Flags().Changed("items-json") && !cmd.Flags().Changed("all-day") && !cmd.Flags().Changed("priority") &&
				!cmd.Flags().Changed("sort-order") {
				return fmt.Errorf("at least one update flag is required")
			}
			payload := ticktick.TaskUpdate{
				ID:        args[0],
				ProjectID: updateProjectID,
			}

			if cmd.Flags().Changed("title") {
				payload.Title = &updateTitle
			}
			if cmd.Flags().Changed("content") {
				payload.Content = &updateContent
			}
			if cmd.Flags().Changed("desc") {
				payload.Desc = &updateDesc
			}
			if cmd.Flags().Changed("start") {
				startDate, err := parseDateTime(updateStart)
				if err != nil {
					return err
				}
				payload.StartDate = &startDate
			}
			if cmd.Flags().Changed("due") {
				dueDate, err := parseDateTime(updateDue)
				if err != nil {
					return err
				}
				payload.DueDate = &dueDate
			}
			if cmd.Flags().Changed("time-zone") {
				payload.TimeZone = &updateTimeZone
			}
			if cmd.Flags().Changed("repeat") {
				payload.RepeatFlag = &updateRepeat
			}
			if cmd.Flags().Changed("kind") {
				kind, err := normalizeTaskKind(updateKind)
				if err != nil {
					return err
				}
				payload.Kind = &kind
			}
			if cmd.Flags().Changed("reminder") {
				payload.Reminders = &updateReminders
			}
			if cmd.Flags().Changed("items-json") {
				items, err := parseChecklistItems(updateItemsJSON)
				if err != nil {
					return err
				}
				payload.Items = &items
			}
			if cmd.Flags().Changed("all-day") {
				payload.IsAllDay = &updateAllDay
			}
			if cmd.Flags().Changed("priority") {
				payload.Priority = &updatePriority
			}
			if cmd.Flags().Changed("sort-order") {
				payload.SortOrder = &updateSortOrder
			}

			updated, err := a.client.UpdateTask(context.Background(), args[0], payload)
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
	update.Flags().StringVar(&updateItemsJSON, "items-json", "", "Checklist items JSON array")
	update.Flags().BoolVar(&updateAllDay, "all-day", false, "Mark task as all-day")
	update.Flags().IntVar(&updatePriority, "priority", 0, "Priority: 0 none, 1 low, 3 medium, 5 high")
	update.Flags().Int64Var(&updateSortOrder, "sort-order", 0, "Task sort order")

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

func parseChecklistItems(raw string) ([]ticktick.ChecklistItem, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var items []ticktick.ChecklistItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, fmt.Errorf("invalid --items-json: %w", err)
	}
	return items, nil
}
