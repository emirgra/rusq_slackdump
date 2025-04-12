package view

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"

	br "github.com/pkg/browser"

	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/bootstrap"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/golang/base"
	"github.com/rusq/slackdump/v3/internal/viewer"
	"github.com/rusq/slackdump/v3/source"
)

//go:embed assets/view.md
var mdView string

var CmdView = &base.Command{
	Short:      "View the slackdump files",
	UsageLine:  "slackdump view [flags]",
	Long:       mdView,
	PrintFlags: true,
	FlagMask:   cfg.OmitAll,
	Run:        runView,
}

var listenAddr string

func init() {
	CmdView.Flag.StringVar(&listenAddr, "listen", "127.0.0.1:8080", "address to listen on")
}

func runView(ctx context.Context, cmd *base.Command, args []string) error {
	if len(args) < 1 {
		base.SetExitStatus(base.SInvalidParameters)
		return fmt.Errorf("viewing slackdump files requires at least one argument")
	}
	flags, err := source.Type(args[0])
	if err != nil {
		base.SetExitStatus(base.SUserError)
		return err
	}

	lg := cfg.Log
	lg.InfoContext(ctx, "opening archive", "source", args[0], "flags", flags)

	src, err := source.Load(ctx, args[0])
	if err != nil {
		base.SetExitStatus(base.SUserError)
		return err
	}
	defer src.Close()

	stoppb := bootstrap.TimedSpinner(ctx, os.Stdout, "Slackdump Viewer is loading files", -1, 0)
	v, err := viewer.New(ctx, listenAddr, src)
	if err != nil {
		base.SetExitStatus(base.SApplicationError)
		return err
	}
	stoppb()
	// sentinel
	go func() {
		<-ctx.Done()
		v.Close()
		lg.InfoContext(ctx, "cleanup complete")
	}()

	lg.InfoContext(ctx, "listening on", "addr", listenAddr)
	go func() {
		if err := br.OpenURL(fmt.Sprintf("http://%s", listenAddr)); err != nil {
			lg.WarnContext(ctx, "unable to open browser", "error", err)
		}
	}()
	if err := v.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			cfg.Log.InfoContext(ctx, "bye")
			return nil
		}
		base.SetExitStatus(base.SApplicationError)
		return err
	}

	return nil
}
