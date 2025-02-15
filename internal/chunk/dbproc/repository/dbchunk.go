package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/rusq/slackdump/v3/internal/chunk"
)

// DBChunk is the database representation of the Chunk.
type DBChunk struct {
	// ID is the unique identifier of the chunk within the session.
	ID         int64           `db:"ID,omitempty"`
	SessionID  int64           `db:"SESSION_ID,omitempty"`
	UnixTS     int64           `db:"UNIX_TS,omitempty"`
	CreatedAt  time.Time       `db:"CREATED_AT,omitempty"`
	TypeID     chunk.ChunkType `db:"TYPE_ID,omitempty"`
	NumRecords int             `db:"NUM_REC"`
	Final      bool            `db:"FINAL"`
}

func (DBChunk) tablename() string {
	return "CHUNK"
}

func (DBChunk) columns() []string {
	return []string{
		"SESSION_ID",
		"UNIX_TS",
		"TYPE_ID",
		"NUM_REC",
		"FINAL",
	}
}

func (d DBChunk) values() []any {
	return []any{
		d.SessionID,
		d.UnixTS,
		d.TypeID,
		d.NumRecords,
		d.Final,
	}
}

type ChunkRepository interface {
	Insert(ctx context.Context, conn sqlx.ExtContext, dbchunk *DBChunk) (int64, error)
}

type chunkRepository struct {
	genericRepository[DBChunk]
}

func NewChunkRepository() ChunkRepository {
	return chunkRepository{newGenericRepository(DBChunk{})}
}

func (r chunkRepository) Insert(ctx context.Context, conn sqlx.ExtContext, dbchunk *DBChunk) (int64, error) {
	stmt := r.stmtInsert()
	res, err := conn.ExecContext(ctx, stmt, dbchunk.values()...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
