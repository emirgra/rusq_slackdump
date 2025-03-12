package convertcmd

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"

	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/bootstrap"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/internal/chunk"
	"github.com/rusq/slackdump/v3/internal/chunk/dbproc"
	"github.com/rusq/slackdump/v3/internal/source"
)

// toDatabase converts the source to the database format.
func toDatabase(ctx context.Context, src, trg string, cflg convertflags) error {
	// detect source type
	st, err := source.Type(src)
	if err != nil {
		return err
	}

	// currently only chunk format is supported for the source.
	if !st.Has(source.FChunk) {
		return ErrSource
	}

	cd, err := chunk.OpenDir(src)
	if err != nil {
		return err
	}
	defer cd.Close()
	dsrc := source.OpenChunkDir(cd, true)
	defer dsrc.Close()

	trg = cfg.StripZipExt(trg)
	if err := chunk2db(ctx, dsrc, trg, cflg); err != nil {
		return err
	}

	if st.Has(source.FMattermost) && cflg.includeFiles {
		slog.Info("Copying files...")
		if err := copyfiles(filepath.Join(trg, chunk.UploadsDir), dsrc.Files().FS()); err != nil {
			return err
		}
	}
	if st.Has(source.FAvatars) && cflg.includeAvatars {
		slog.Info("Copying avatars...")
		if err := copyfiles(filepath.Join(trg, chunk.AvatarsDir), dsrc.Avatars().FS()); err != nil {
			return err
		}
	}
	return nil
}

// chunk2db converts the chunk source to the database format, it creates the
// database in the directory dir.
func chunk2db(ctx context.Context, src *source.ChunkDir, dir string, cflg convertflags) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	remove := true
	defer func() {
		// remove on failed conversion
		if remove {
			_ = os.RemoveAll(dir)
		}
	}()

	slog.Info("output", "database", filepath.Join(dir, "slackdump.sqlite"))

	// create a new database
	wconn, si, err := bootstrap.Database(dir, "convert")
	if err != nil {
		return err
	}
	defer wconn.Close()

	dbp, err := dbproc.New(ctx, wconn, si)
	if err != nil {
		return err
	}
	defer dbp.Close()

	txx, err := wconn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer txx.Rollback()

	enc := &encoder{dbp: dbp, tx: txx}
	if err := src.ToChunk(ctx, enc, cflg.sessionID); err != nil {
		return err
	}
	if err := txx.Commit(); err != nil {
		return err
	}

	remove = false
	return nil
}

// encoder implements the chunk.Encoder around the unsafe database insert.
// It operates in a single transaction tx.
type encoder struct {
	dbp *dbproc.DBP
	tx  *sqlx.Tx
}

func (e *encoder) Encode(ctx context.Context, ch *chunk.Chunk) error {
	_, err := e.dbp.UnsafeInsertChunk(ctx, e.tx, ch)
	return err
}
