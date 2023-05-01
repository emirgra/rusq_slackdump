package dirproc

import (
	"io"
	"sync/atomic"

	"github.com/rusq/slackdump/v2/internal/chunk"
)

// baseproc exposes recording functionality for processor, and handles chunk
// file creation.
type baseproc struct {
	dir    string
	wc     io.WriteCloser
	closed atomic.Bool
	*chunk.Recorder
}

// newBaseProc initialises the new base processor.  It creates a new chunk file
// in a directory dir which must exist.
func newBaseProc(dir string, name chunk.FileID) (*baseproc, error) {
	cd, err := chunk.OpenDir(dir)
	if err != nil {
		return nil, err
	}
	wc, err := cd.Create(name)
	if err != nil {
		return nil, err
	}

	r := chunk.NewRecorder(wc)
	return &baseproc{dir: dir, wc: wc, Recorder: r}, nil
}

// Close closes the processor and the underlying chunk file.
func (p *baseproc) Close() error {
	if p.closed.Load() {
		return nil
	}
	if err := p.Recorder.Close(); err != nil {
		p.wc.Close()
		return err
	}
	p.closed.Store(true)
	if err := p.wc.Close(); err != nil {
		return err
	}
	return nil
}