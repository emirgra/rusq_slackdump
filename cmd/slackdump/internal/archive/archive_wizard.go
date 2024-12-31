package archive

import (
	"context"

	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/golang/base"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/ui/cfgui"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/ui/dumpui"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/ui/updaters"
	"github.com/rusq/slackdump/v3/internal/structures"
)

func archiveWizard(ctx context.Context, cmd *base.Command, args []string) error {
	w := &dumpui.Wizard{
		Title:       "Archive Slack Workspace",
		Name:        "Archive",
		Cmd:         cmd,
		LocalConfig: configuration,
		ArgsFn: func() []string {
			if len(entryList) > 0 {
				return structures.SplitEntryList(entryList)
			}
			return nil
		},
	}
	return w.Run(ctx)
}

var entryList string

func configuration() cfgui.Configuration {
	return cfgui.Configuration{
		cfgui.ParamGroup{
			Name: "Optional parameters",
			Params: []cfgui.Parameter{
				cfgui.ChannelIDs(&entryList, false),
				{
					Name:        "Member Only",
					Value:       cfgui.Checkbox(cfg.MemberOnly),
					Description: "Export only channels, which current user belongs to",
					Inline:      true,
					Updater:     updaters.NewBool(&cfg.MemberOnly),
				},
			},
		},
	}
}