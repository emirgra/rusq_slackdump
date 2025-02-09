package repository

import (
	"encoding/json"

	"github.com/rusq/slack"
)

type DBSearchMessage struct {
	ID          int64   `db:"ID"`
	ChunkID     int64   `db:"CHUNK_ID"`
	LoadDTTM    string  `db:"LOAD_DTTM,omitempty"`
	ChannelID   string  `db:"CHANNEL_ID"`
	ChannelName *string `db:"CHANNEL_NAME,omitempty"`
	TS          string  `db:"TS"`
	Text        *string `db:"TXT,omitempty"`
	IDX         int     `db:"IDX"`
	Data        []byte  `db:"DATA"`
}

func NewDBSearchMessage(chunkID int64, idx int, sm *slack.SearchMessage) (*DBSearchMessage, error) {
	data, err := json.Marshal(sm)
	if err != nil {
		return nil, err
	}
	return &DBSearchMessage{
		ChunkID:     chunkID,
		ChannelID:   sm.Channel.ID,
		ChannelName: orNull(sm.Channel.Name != "", sm.Channel.Name),
		TS:          sm.Timestamp,
		Text:        orNull(sm.Text != "", sm.Text),
		IDX:         idx,
		Data:        data,
	}, nil
}

func (DBSearchMessage) Table() string {
	return "SEARCH_MESSAGE"
}

func (DBSearchMessage) Columns() []string {
	return []string{"CHUNK_ID", "CHANNEL_ID", "CHANNEL_NAME", "TS", "TXT", "IDX", "DATA"}
}

func (c DBSearchMessage) Values() []any {
	return []interface{}{c.ChunkID, c.ChannelID, c.ChannelName, c.TS, c.Text, c.IDX, c.Data}
}

type SearchMessageRepository interface {
	repository[*DBSearchMessage]
}

func NewSearchMessageRepository() SearchMessageRepository {
	return newGenericRepository[*DBSearchMessage]()
}
