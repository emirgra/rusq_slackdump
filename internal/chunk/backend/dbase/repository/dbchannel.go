package repository

import (
	"context"
	"iter"

	"github.com/jmoiron/sqlx"
	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3/internal/chunk"
)

type DBChannel struct {
	ID      string  `db:"ID"`
	ChunkID int64   `db:"CHUNK_ID"`
	Name    *string `db:"NAME"`
	Index   int     `db:"IDX"`
	Data    []byte  `db:"DATA"`
}

func NewDBChannel(chunkID int64, n int, channel *slack.Channel) (*DBChannel, error) {
	data, err := marshal(channel)
	if err != nil {
		return nil, err
	}
	return &DBChannel{
		ID:      channel.ID,
		ChunkID: chunkID,
		Name:    orNull(channel.Name != "", channel.Name),
		Index:   n,
		Data:    data,
	}, nil
}

func (c DBChannel) tablename() string {
	return "CHANNEL"
}

func (c DBChannel) userkey() []string {
	return slice("ID")
}

func (c DBChannel) columns() []string {
	return []string{"ID", "CHUNK_ID", "NAME", "IDX", "DATA"}
}

func (c DBChannel) values() []any {
	return []any{c.ID, c.ChunkID, c.Name, c.Index, c.Data}
}

func (c DBChannel) Val() (slack.Channel, error) {
	return unmarshalt[slack.Channel](c.Data)
}

//go:generate mockgen -destination=mock_repository/mock_channel.go . ChannelRepository
type ChannelRepository interface {
	BulkRepository[DBChannel]
}

type channelRepository struct {
	genericRepository[DBChannel]
}

func NewChannelRepository() ChannelRepository {
	return channelRepository{newGenericRepository(DBChannel{})}
}

func (r channelRepository) Count(ctx context.Context, conn sqlx.QueryerContext) (int64, error) {
	return r.CountType(ctx, conn, chunk.CChannelInfo)
}

func (r channelRepository) Get(ctx context.Context, conn sqlx.ExtContext, id any) (DBChannel, error) {
	return r.GetType(ctx, conn, id, chunk.CChannelInfo)
}

func (r channelRepository) AllOfType(ctx context.Context, conn sqlx.QueryerContext, typeID ...chunk.ChunkType) (iter.Seq2[DBChannel, error], error) {
	return r.allOfTypeWhere(ctx, conn, queryParams{OrderBy: []string{"T.NAME"}}, typeID...)
}
