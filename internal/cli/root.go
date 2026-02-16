package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/apktdev/ticktick-cli/internal/config"
	"github.com/apktdev/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

type app struct {
	cfg           *config.Config
	client        *ticktick.Client
	jsonMode      bool
	shouldPersist bool
}

func NewRootCmd() *cobra.Command {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	envOverride, err := config.ApplyEnvOverrides(cfg)
	if err != nil {
		panic(err)
	}

	a := &app{
		cfg:           cfg,
		client:        ticktick.New(cfg),
		shouldPersist: !envOverride,
	}
	root := &cobra.Command{
		Use:   "ticktick",
		Short: "TickTick command line client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !a.shouldPersist {
				return nil
			}
			return config.Save(a.cfg)
		},
	}
	root.PersistentFlags().BoolVar(&a.jsonMode, "json", false, "Output JSON")

	root.AddCommand(a.newAuthCmd())
	root.AddCommand(a.newProjectsCmd())
	root.AddCommand(a.newTasksCmd())
	return root
}

func (a *app) print(v any) error {
	if a.jsonMode {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}
	switch x := v.(type) {
	case string:
		fmt.Println(x)
	case []ticktick.Project:
		for _, p := range x {
			fmt.Printf("%s\t%s\n", p.ID, p.Name)
		}
	case []ticktick.Task:
		for _, t := range x {
			due := ""
			if t.DueDate != "" {
				due = " due=" + t.DueDate
			}
			fmt.Printf("%s\t%s%s\n", t.ID, t.Title, due)
		}
	default:
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
	}
	return nil
}

func parseDateTime(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, l := range layouts {
		t, err := time.Parse(l, input)
		if err == nil {
			return t.Format("2006-01-02T15:04:05-0700"), nil
		}
	}
	return "", fmt.Errorf("invalid date format %q; use RFC3339 or YYYY-MM-DD", input)
}

func required(val, name string) error {
	if val == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
