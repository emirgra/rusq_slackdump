package dirproc

import (
	"io"
	"sync/atomic"

	"github.com/rusq/slackdump/v2/internal/chunk"
)

// baseproc exposes recording functionality for processor, and handles chunk
// file creation.
type baseproc struct {
	// cd     *chunk.Directory
	wc     io.WriteCloser
	closed atomic.Bool
	*chunk.Recorder
}

// newBaseProc initialises the new base processor.  It creates a new chunk file
// in a directory dir which must exist.
func newBaseProc(cd *chunk.Directory, name chunk.FileID) (*baseproc, error) {
	wc, err := cd.Create(name)
	if err != nil {
		return nil, err
	}

	r := chunk.NewRecorder(wc)
	return &baseproc{
		// cd: cd,
		wc:       wc,
		Recorder: r,
	}, nil
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
