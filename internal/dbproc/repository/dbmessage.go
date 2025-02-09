package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3/internal/fasttime"
	"github.com/rusq/slackdump/v3/internal/structures"
)

type DBMessage struct {
	ID        int64     `db:"ID,omitempty"`
	ChunkID   int64     `db:"CHUNK_ID,omitempty"`
	LoadDTTM  time.Time `db:"LOAD_DTTM,omitempty"`
	ChannelID string    `db:"CHANNEL_ID"`
	TS        string    `db:"TS"`
	ParentID  *int64    `db:"PARENT_ID,omitempty"`
	ThreadTS  *string   `db:"THREAD_TS,omitempty"`
	IsParent  bool      `db:"IS_PARENT"`
	Index     int       `db:"IDX"`
	NumFiles  int       `db:"NUM_FILES"`
	Text      string    `db:"TXT"`
	Data      []byte    `db:"DATA"`
}

func NewDBMessage(dbchunkID int64, idx int, channelID string, msg *slack.Message) (*DBMessage, error) {
	ts, err := fasttime.TS2int(msg.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("insertMessages fasttime: %w", err)
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("insertMessages marshal: %w", err)
	}
	var (
		parentID *int64
	)
	if msg.ThreadTimestamp != "" {
		if parID, err := fasttime.TS2int(msg.ThreadTimestamp); err != nil {
			return nil, fmt.Errorf("insertMessages fasttime thread: %w", err)
		} else {
			parentID = &parID
		}
	}

	dbm := DBMessage{
		ID:        ts,
		ChunkID:   dbchunkID,
		ChannelID: channelID,
		TS:        msg.Timestamp,
		ParentID:  parentID,
		ThreadTS:  orNull(msg.ThreadTimestamp != "", msg.ThreadTimestamp),
		IsParent:  structures.IsThreadStart(msg),
		Index:     idx,
		NumFiles:  len(msg.Files),
		Text:      msg.Text,
		Data:      data,
	}
	return &dbm, nil
}

func (dbm *DBMessage) Table() string {
	return "MESSAGE"
}

func (dbm *DBMessage) Columns() []string {
	return []string{
		"ID",
		"CHUNK_ID",
		"CHANNEL_ID",
		"TS",
		"PARENT_ID",
		"THREAD_TS",
		"IS_PARENT",
		"IDX",
		"NUM_FILES",
		"TXT",
		"DATA",
	}
}

func (dbm *DBMessage) Values() []any {
	return []any{
		dbm.ID,
		dbm.ChunkID,
		dbm.ChannelID,
		dbm.TS,
		dbm.ParentID,
		dbm.ThreadTS,
		dbm.IsParent,
		dbm.Index,
		dbm.NumFiles,
		dbm.Text,
		dbm.Data,
	}
}

type MessageRepository interface {
	Insert(ctx context.Context, conn sqlx.ExtContext, m *DBMessage) error
	InsertAll(ctx context.Context, tx PrepareExtContext, mm iter.Seq2[*DBMessage, error]) (int, error)
}

type messageRepository struct {
	repository[*DBMessage]
}

func NewMessageRepository() MessageRepository {
	return messageRepository{newGenericRepository[*DBMessage]()}
}
