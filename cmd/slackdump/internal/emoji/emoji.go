package emoji

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/rusq/fsadapter"
	"github.com/schollz/progressbar/v3"

	"github.com/rusq/slackdump/v3"
	"github.com/rusq/slackdump/v3/auth"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/bootstrap"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/emoji/emojidl"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/golang/base"
	"github.com/rusq/slackdump/v3/internal/edge"
)

//go:embed assets/emoji.md
var emojiMD string

var CmdEmoji = &base.Command{
	Run:         run,
	UsageLine:   "slackdump emoji [flags]",
	Short:       "download workspace emoticons ಠ_ಠ",
	Long:        emojiMD,
	FlagMask:    cfg.OmitAll &^ cfg.OmitAuthFlags &^ cfg.OmitOutputFlag &^ cfg.OmitWorkspaceFlag &^ cfg.OmitWithFilesFlag,
	RequireAuth: true,
	PrintFlags:  true,
}

type options struct {
	full bool
	emojidl.Options
}

// emoji specific flags
var cmdFlags = options{
	Options: emojidl.Options{
		FailFast: false,
	},
}

func init() {
	CmdEmoji.Wizard = wizard
	CmdEmoji.Flag.BoolVar(&cmdFlags.FailFast, "ignore-errors", true, "ignore download errors (skip failed emojis)")
	CmdEmoji.Flag.BoolVar(&cmdFlags.full, "full", false, "fetch emojis using Edge API to get full emoji information, including usernames")
}

func run(ctx context.Context, cmd *base.Command, args []string) error {
	fsa, err := fsadapter.New(cfg.Output)
	if err != nil {
		return err
	}
	defer fsa.Close()

	prov, err := auth.FromContext(ctx)
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}

	cmdFlags.WithFiles = cfg.WithFiles

	start := time.Now()
	r, cl := statusReporter()
	defer cl.Close()
	if cmdFlags.full {
		err = runEdge(ctx, fsa, prov, r)
	} else {
		err = runLegacy(ctx, fsa, r)
	}
	cl.Close()
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}

	slog.InfoContext(ctx, "Emojis downloaded", "dir", cfg.Output, "took", time.Since(start).String())
	return nil
}

func statusReporter() (emojidl.StatusFunc, io.Closer) {
	pb := progressbar.NewOptions(0,
		progressbar.OptionSetDescription("Downloading emojis"),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowCount(),
	)
	var once sync.Once
	return func(name string, total, count int) {
		once.Do(func() {
			pb.ChangeMax(total)
		})
		pb.Add(1)
	}, pb
}

func runLegacy(ctx context.Context, fsa fsadapter.FS, cb emojidl.StatusFunc) error {
	sess, err := bootstrap.SlackdumpSession(ctx, slackdump.WithFilesystem(fsa))
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}

	return emojidl.DlFS(ctx, sess, fsa, &cmdFlags.Options, cb)
}

func runEdge(ctx context.Context, fsa fsadapter.FS, prov auth.Provider, cb emojidl.StatusFunc) error {
	sess, err := edge.New(ctx, prov)
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}
	defer sess.Close()

	if err := emojidl.DlEdgeFS(ctx, sess, fsa, &cmdFlags.Options, cb); err != nil {
		base.SetExitStatus(base.SApplicationError)
		return fmt.Errorf("application error: %s", err)
	}
	return nil
}
