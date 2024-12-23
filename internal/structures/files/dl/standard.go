package dl

// standard file export

import (
	"errors"
	"log/slog"
	"path"
	"path/filepath"

	"github.com/rusq/slack"

	"github.com/rusq/fsadapter"
	"github.com/rusq/slackdump/v3"
	"github.com/rusq/slackdump/v3/downloader"
	"github.com/rusq/slackdump/v3/internal/structures/files"
	"github.com/rusq/slackdump/v3/types"
)

type Std struct {
	base
}

// NewStd returns standard dl, which downloads files into
// "channel_id/attachments" directory.
func NewStd(fs fsadapter.FS, cl *slack.Client, l *slog.Logger, token string) *Std {
	return &Std{
		base: base{
			dl:    downloader.NewV1(cl, fs, downloader.LoggerV1(l)),
			l:     l,
			token: token,
		},
	}
}

// ProcessFunc returns the function that downloads the file into
// channel_id/attachments directory. If Slack token is set, it updates the
// thumbnails to include that token.  It replaces the file URL to point to
// physical downloaded files on disk.
func (d *Std) ProcessFunc(channelName string) slackdump.ProcessFunc {
	const (
		dirAttach = "attachments"
	)

	dir := filepath.Join(channelName, dirAttach)
	return func(msg []types.Message, _ string) (slackdump.ProcessResult, error) {
		total := 0
		if err := files.Extract(msg, files.Root, func(file slack.File, addr files.Addr) error {
			filename, err := d.dl.DownloadFile(dir, file)
			if err != nil {
				return err
			}
			d.l.Debug("submitted for download", "filename", file.Name)
			total++
			if d.token != "" {
				if err := files.Update(msg, addr, files.UpdateTokenFn(d.token)); err != nil {
					return err
				}
			}
			return files.Update(msg, addr, files.UpdatePathFn(path.Join(dirAttach, path.Base(filename))))
		}); err != nil {
			if errors.Is(err, downloader.ErrNotStarted) {
				return slackdump.ProcessResult{Entity: entFiles, Count: 0}, nil
			}
			return slackdump.ProcessResult{}, err
		}

		return slackdump.ProcessResult{Entity: entFiles, Count: total}, nil
	}
}
