// Package fileproc is the file processor that can be used in conjunction with
// the transformer.  It downloads files to the local filesystem using the
// provided downloader.  Probably it's a good idea to use the
// [downloader.Client] for this.
package fileproc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/rusq/fsadapter"
	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3/downloader"
	"github.com/rusq/slackdump/v3/internal/structures/files"
)

// Downloader is the interface that wraps the Download method.
type Downloader interface {
	// Download should download the file at the given URL and save it to the
	// given path.
	Download(fullpath string, url string) error
}

// Subprocessor is the file subprocessor, that downloads files to the path
// returned by the filepath function.
// Zero value of this type is not usable.
type Subprocessor struct {
	dcl      Downloader
	filepath func(ci *slack.Channel, f *slack.File) string
}

// NewSubprocessor initialises the subprocessor.
func NewSubprocessor(dl Downloader, fp func(ci *slack.Channel, f *slack.File) string) Subprocessor {
	if fp == nil {
		panic("filepath function is nil")
	}
	return Subprocessor{
		dcl:      dl,
		filepath: fp,
	}
}

func (b Subprocessor) Files(ctx context.Context, channel *slack.Channel, msg slack.Message, ff []slack.File) error {
	for _, f := range ff {
		if !IsValid(&f) {
			continue
		}
		if err := b.dcl.Download(b.filepath(channel, &f), f.URLPrivateDownload); err != nil {
			return err
		}
	}
	return nil
}

// PathUpdateFunc updates the path in URLDownload and URLPrivateDownload of every
// file in the given message slice to point to the physical downloaded file
// location.  It can be plugged in the pipeline of Dump.
func (b Subprocessor) PathUpdateFunc(channelID, threadTS string, mm []slack.Message) error {
	for i := range mm {
		for j := range mm[i].Files {
			ch := new(slack.Channel)
			ch.ID = channelID
			path := b.filepath(ch, &mm[i].Files[j])
			if err := files.UpdatePathFn(path)(&mm[i].Files[j]); err != nil {
				return err
			}
		}
	}
	return nil
}

// ExportTokenUpdateFn returns a function that appends the token to every file
// URL in the given message.
func ExportTokenUpdateFn(token string) func(msg *slack.Message) error {
	fn := files.UpdateTokenFn(token)
	return func(msg *slack.Message) error {
		for i := range msg.Files {
			if err := fn(&msg.Files[i]); err != nil {
				return err
			}
		}
		return nil
	}
}

var invalidModes = map[string]struct{}{
	"hidden_by_limit": {},
	"external":        {},
	"tombstone":       {},
}

// IsValid returns true if the file can be downloaded and is valid.
func IsValid(f *slack.File) bool {
	return IsValidWithReason(f) == nil
}

func IsValidWithReason(f *slack.File) error {
	if f == nil {
		return errors.New("file is nil")
	}
	if _, ok := invalidModes[f.Mode]; ok {
		return fmt.Errorf("invalid file mode %q", f.Mode)
	}
	if !f.IsExternal && f.Name == "" {
		return fmt.Errorf("invalid file: external=%v, name=%q", f.IsExternal, f.Name)
	}
	return nil
}

type NoopDownloader struct{}

func (NoopDownloader) Download(fullpath string, url string) error {
	return nil
}

type FileGetter interface {
	// GetFile retreives a given file from its private download URL
	GetFileContext(ctx context.Context, downloadURL string, writer io.Writer) error
}

// NewDownloader initializes the downloader and returns it, along with a
// function that should be called to stop it.
func NewDownloader(ctx context.Context, gEnabled bool, cl FileGetter, fsa fsadapter.FS, lg *slog.Logger) (sdl Downloader, stop func()) {
	if !gEnabled {
		return NoopDownloader{}, func() {}
	} else {
		dl := downloader.New(cl, fsa, downloader.WithLogger(lg))
		if err := dl.Start(ctx); err != nil {
			lg.Error("failed to start downloader", "error", err)
			return NoopDownloader{}, func() {}
		}
		return dl, dl.Stop
	}
}
