package repository

import (
	"encoding/json"
	"time"

	"github.com/rusq/slack"
)

type DBWorkspace struct {
	ID           int64     `db:"ID,omitempty"`
	ChunkID      int64     `db:"CHUNK_ID"`
	LoadDTTM     time.Time `db:"LOAD_DTTM,omitempty"`
	Team         string    `db:"TEAM"`
	User         string    `db:"USERNAME"`
	TeamID       string    `db:"TEAM_ID"`
	UserID       string    `db:"USER_ID"`
	EnterpriseID *string   `db:"ENTERPRISE_ID"`
	URL          string    `db:"URL"`
	Data         []byte    `db:"DATA"`
}

func NewDBWorkspace(chunkID int64, wi *slack.AuthTestResponse) (*DBWorkspace, error) {
	data, err := json.Marshal(wi)
	if err != nil {
		return nil, err
	}
	return &DBWorkspace{
		ChunkID:      chunkID,
		Team:         wi.Team,
		User:         wi.UserID,
		TeamID:       wi.TeamID,
		UserID:       wi.UserID,
		EnterpriseID: orNull(wi.EnterpriseID != "", wi.EnterpriseID),
		URL:          wi.URL,
		Data:         data,
	}, nil
}

func (w DBWorkspace) Table() string {
	return "WORKSPACE"
}

func (w DBWorkspace) Columns() []string {
	return []string{
		"CHUNK_ID",
		"TEAM",
		"USER_ID",
		"TEAM_ID",
		"USER_ID",
		"ENTERPRISE_ID",
		"URL",
		"DATA",
	}
}

func (w DBWorkspace) Values() []any {
	return []any{
		w.ChunkID,
		w.Team,
		w.UserID,
		w.TeamID,
		w.UserID,
		w.EnterpriseID,
		w.URL,
		w.Data,
	}
}

type WorkspaceRepository interface {
	repository[*DBWorkspace]
}

func NewWorkspaceRepository() WorkspaceRepository {
	return newGenericRepository[*DBWorkspace]()
}
