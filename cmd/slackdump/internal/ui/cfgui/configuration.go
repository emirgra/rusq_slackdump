package cfgui

import (
	"context"
	"errors"
	"os"
	"runtime/trace"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/apiconfig"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/ui/bubbles/filemgr"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/ui/updaters"
)

type Configuration []ParamGroup

type ParamGroup struct {
	Name   string
	Params []Parameter
}

type Parameter struct {
	Name        string
	Value       string
	Description string
	Inline      bool
	Updater     tea.Model
}

func globalConfig() Configuration {
	cwd, _ := os.Getwd()
	return Configuration{
		{
			Name: "Timeframe",
			Params: []Parameter{
				{
					Name:        "Start date",
					Value:       cfg.Oldest.String(),
					Description: "The oldest message to fetch",
					Updater:     updaters.NewDTTM((*time.Time)(&cfg.Oldest)),
				},
				{
					Name:        "End date",
					Value:       cfg.Latest.String(),
					Description: "The newest message to fetch",
					Updater:     updaters.NewDTTM((*time.Time)(&cfg.Latest)),
				},
			},
		},
		{
			Name: "Output",
			Params: []Parameter{
				{
					Name:        "Output",
					Value:       cfg.Output,
					Inline:      true,
					Description: "Output directory",
					Updater:     updaters.NewFileNew(&cfg.Output, "ZIP or Directory", false, true),
				},
			},
		},
		{
			Name: "API options",
			Params: []Parameter{
				{
					Name:        "Download files",
					Value:       Checkbox(cfg.WithFiles),
					Description: "Download files",
					Updater:     updaters.NewBool(&cfg.WithFiles),
				},
				{
					Name:        "Enterprise mode",
					Value:       Checkbox(cfg.ForceEnterprise),
					Description: "Force enterprise mode",
					Updater:     updaters.NewBool(&cfg.ForceEnterprise),
				},
				{
					Name:        "API limits file",
					Value:       cfg.ConfigFile,
					Description: "API limits file",
					Updater: updaters.NewFilepickModel(
						&cfg.ConfigFile,
						filemgr.New(os.DirFS(cwd), cwd, ".", 15, apiconfig.ConfigExts...),
						validateAPIconfig,
					),
				},
			},
		},
		{
			Name: "Cache Control",
			Params: []Parameter{
				{
					Name:        "Local Cache Directory",
					Value:       cfg.LocalCacheDir,
					Description: "Local Cache Directory for user data",
				},
				{
					Name:        "User Cache Retention",
					Value:       cfg.UserCacheRetention.String(),
					Description: "For how long user cache is kept, until it is fetched again",
					Inline:      true,
					Updater:     updaters.NewDuration(&cfg.UserCacheRetention, false),
				},
				{
					Name:        "Disable User Cache",
					Value:       Checkbox(cfg.NoUserCache),
					Description: "Disable User Cache",
					Updater:     updaters.NewBool(&cfg.NoUserCache),
				},
				{
					Name:        "Disable Chunk Cache",
					Value:       Checkbox(cfg.NoChunkCache),
					Description: "Disable Chunk Cache",
					Updater:     updaters.NewBool(&cfg.NoChunkCache),
				},
			},
		},
	}
}

func validateAPIconfig(s string) error {
	_, task := trace.NewTask(context.Background(), "validateAPIconfig")
	defer task.End()
	if s == "" {
		return nil
	}
	if _, err := os.Stat(s); err != nil {
		return err
	}
	if err := apiconfig.CheckFile(s); err != nil {
		return errors.New("not a valid API limits configuration file")
	}
	return nil
}
