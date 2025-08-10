package cfgui

import (
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/ui/updaters"
	"github.com/rusq/slackdump/v3/internal/structures"
)

// Common reusable parameters

func ChannelIDs(v *string, required bool) Parameter {
	name := "Channel IDs or URLs"
	descr := "List of channel IDs or URLs to dump"
	if required {
		name = "* " + name
		descr = descr + " (REQUIRED)"
	}
	return Parameter{
		Name:        name,
		Value:       *v,
		Description: descr,
		Inline:      true,
		Updater:     updaters.NewString(v, "", false, structures.ValidateEntityList),
	}
}

// MemberOnly returns a checkbox parameter for Member Only flag.
func MemberOnly() Parameter {
	return Parameter{
		Name:        "Member Only",
		Value:       Checkbox(cfg.MemberOnly),
		Description: "Export only channels, which you belongs to.",
		Updater:     updaters.NewBool(&cfg.MemberOnly),
	}
}

// RecordFiles returns a checkbox parameter for Record Files flag.
func RecordFiles() Parameter {
	return Parameter{
		Name:        "Record Files",
		Value:       Checkbox(cfg.RecordFiles),
		Description: "Record file chunks in chunk files.",
		Updater:     updaters.NewBool(&cfg.RecordFiles),
	}
}

func Avatars() Parameter {
	return Parameter{
		Name:        "Download Avatars",
		Value:       Checkbox(cfg.WithAvatars),
		Description: "Download avatars.",
		Updater:     updaters.NewBool(&cfg.WithAvatars),
	}
}

func OnlyChannelUsers() Parameter {
	return Parameter{
		Name:        "Only Channel Users",
		Value:       Checkbox(cfg.OnlyChannelUsers),
		Description: "Only users participating in visible conversations are exported.",
		Updater:     updaters.NewBool(&cfg.OnlyChannelUsers),
	}
}
