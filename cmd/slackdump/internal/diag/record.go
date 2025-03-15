package diag

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rusq/slackdump/v3/internal/chunk/backend/dbase"
	"github.com/rusq/slackdump/v3/internal/chunk/backend/dbase/repository"

	"github.com/jmoiron/sqlx"

	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/bootstrap"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/golang/base"
	"github.com/rusq/slackdump/v3/internal/chunk"
	"github.com/rusq/slackdump/v3/internal/structures"
)

var cmdRecord = &base.Command{
	UsageLine:  "slackdump tools record",
	Short:      "chunk record commands",
	Commands:   []*base.Command{cmdRecordStream},
	HideWizard: true,
}

var cmdRecordStream = &base.Command{
	UsageLine: "slackdump tools record stream [options] <channel>",
	Short:     "dump slack data in a chunk record format",
	Long: `
# Record tool

Records the data from a channel in a chunk record format.

See also: slackdump tool obfuscate
`,
	FlagMask:    cfg.OmitOutputFlag | cfg.OmitDownloadFlag,
	PrintFlags:  true,
	RequireAuth: true,
}

func init() {
	// break init cycle
	cmdRecordStream.Run = runRecord
}

var output = cmdRecordStream.Flag.String("output", "", "output file")

func runRecord(ctx context.Context, _ *base.Command, args []string) error {
	if len(args) == 0 {
		base.SetExitStatus(base.SInvalidParameters)
		return errors.New("missing channel argument")
	}

	sess, err := bootstrap.SlackdumpSession(ctx)
	if err != nil {
		base.SetExitStatus(base.SInitializationError)
		return err
	}

	// var w io.Writer
	// if *output == "" {
	// 	w = os.Stdout
	// } else {
	// 	if f, err := os.Create(*output); err != nil {
	// 		base.SetExitStatus(base.SApplicationError)
	// 		return err
	// 	} else {
	// 		defer f.Close()
	// 		w = f
	// 	}
	// }

	db, err := sqlx.Open(repository.Driver, "record.db")
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}
	defer db.Close()

	runParams := dbase.SessionInfo{
		FromTS:         (*time.Time)(&cfg.Oldest),
		ToTS:           (*time.Time)(&cfg.Latest),
		FilesEnabled:   cfg.WithFiles,
		AvatarsEnabled: cfg.WithAvatars,
		Mode:           "record",
		Args:           strings.Join(os.Args, "|"),
	}

	p, err := dbase.New(ctx, db, runParams)
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}
	defer p.Close()

	// rec := chunk.NewRecorder(w)
	rec := chunk.NewCustomRecorder(p)
	for _, ch := range args {
		lg := cfg.Log.With("channel_id", ch)
		lg.InfoContext(ctx, "streaming")
		if err := sess.Stream().SyncConversations(ctx, rec, structures.EntityItem{Id: ch}); err != nil {
			if err2 := rec.Close(); err2 != nil {
				base.SetExitStatus(base.SApplicationError)
				return fmt.Errorf("error streaming channel %q: %w; error closing recorder: %v", ch, err, err2)
			}
			return err
		}
	}
	if err := rec.Close(); err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}
	return nil
}
