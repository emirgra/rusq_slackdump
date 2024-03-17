package source

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/rusq/slack"
	"github.com/rusq/slackdump/v3/export"
	"github.com/rusq/slackdump/v3/internal/structures"
)

// Export implements viewer.Sourcer for the zip file Slack export format.
type Export struct {
	fs        fs.FS
	channels  []slack.Channel
	chanNames map[string]string // maps the channel id to the channel name.
	name      string            // name of the file
	idx       structures.ExportIndex
	filestorage
}

func NewExport(fsys fs.FS, name string) (*Export, error) {
	var idx structures.ExportIndex
	if err := idx.Unmarshal(fsys); err != nil {
		return nil, err
	}
	chans := idx.Restore()
	z := &Export{
		fs:        fsys,
		name:      name,
		idx:       idx,
		channels:  chans,
		chanNames: make(map[string]string, len(chans)),
	}
	// initialise channels for quick lookup
	for _, ch := range z.channels {
		z.chanNames[ch.ID] = structures.NVL(ch.Name, ch.ID)
	}
	// determine files path
	fst, err := loadStorage(fsys)
	if err != nil {
		return nil, err
	}
	z.filestorage = fst

	return z, nil
}

// loadStorage determines the type of the file storage used and initialises
// appropriate filestorage implementation.
func loadStorage(fsys fs.FS) (filestorage, error) {
	if _, err := fs.Stat(fsys, "__uploads"); err == nil {
		return newMattermostStorage(fsys)
	}
	idx, err := buildFileIndex(fsys, ".")
	if err != nil || len(idx) == 0 {
		return fstNotFound{}, nil
	}
	return newStandardStorage(fsys, idx), nil
}

func (e *Export) Channels() ([]slack.Channel, error) {
	return e.channels, nil
}

func (e *Export) Users() ([]slack.User, error) {
	return e.idx.Users, nil
}

func (e *Export) Close() error {
	return nil
}

func (e *Export) Name() string {
	return e.name
}

func (e *Export) Type() string {
	return "export"
}

func (e *Export) AllMessages(channelID string) ([]slack.Message, error) {
	// find the channel
	name, ok := e.chanNames[channelID]
	if !ok {
		return nil, fmt.Errorf("%w: %s", fs.ErrNotExist, channelID)
	}
	_, err := fs.Stat(e.fs, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", fs.ErrNotExist, name)
	}

	var mm []slack.Message
	if err := fs.WalkDir(e.fs, name, func(pth string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || path.Ext(pth) != ".json" {
			return nil
		}
		// read the file
		em, err := unmarshal[[]export.ExportMessage](e.fs, pth)
		if err != nil {
			return err
		}
		for _, m := range em {
			mm = append(mm, slack.Message{Msg: *m.Msg})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("AllMessages: walk: %s", err)
	}
	return mm, nil
}

func (e *Export) AllThreadMessages(channelID, threadID string) ([]slack.Message, error) {
	m, err := e.AllMessages(channelID)
	if err != nil {
		return nil, err
	}
	var tm []slack.Message
	for _, msg := range m {
		if msg.ThreadTimestamp == threadID {
			tm = append(tm, msg)
		}
	}
	return tm, nil
}

func (e *Export) ChannelInfo(channelID string) (*slack.Channel, error) {
	c, err := e.Channels()
	if err != nil {
		return nil, err
	}
	for _, ch := range c {
		if ch.ID == channelID {
			return &ch, nil
		}
	}
	return nil, fmt.Errorf("%s: %s", "channel not found", channelID)
}