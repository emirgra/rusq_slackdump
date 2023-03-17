package slackdump

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"errors"

	"github.com/golang/mock/gomock"
	"github.com/slack-go/slack"

	"github.com/rusq/slackdump/v2/internal/cache"
	"github.com/rusq/slackdump/v2/internal/fixtures"
	"github.com/rusq/slackdump/v2/internal/structures"
	"github.com/rusq/slackdump/v2/logger"
	"github.com/rusq/slackdump/v2/types"
)

const testSuffix = "UNIT"

var testUsers = types.Users(fixtures.TestUsers)

func TestUsers_IndexByID(t *testing.T) {
	users := []slack.User{
		{ID: "USLACKBOT", Name: "slackbot"},
		{ID: "USER2", Name: "User 2"},
	}
	tests := []struct {
		name string
		us   types.Users
		want structures.UserIndex
	}{
		{"test 1", users, structures.UserIndex{
			"USLACKBOT": &users[0],
			"USER2":     &users[1],
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.IndexByID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.MakeUserIDIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_fetchUsers(t *testing.T) {
	type fields struct {
		Users   types.Users
		options Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expectFn func(*mockClienter)
		want     types.Users
		wantErr  bool
	}{
		{
			"ok",
			fields{options: DefOptions},
			args{context.Background()},
			func(mc *mockClienter) {
				mc.EXPECT().GetUsersContext(gomock.Any()).Return([]slack.User(testUsers), nil)
			},
			testUsers,
			false,
		},
		{
			"api error",
			fields{options: DefOptions},
			args{context.Background()},
			func(mc *mockClienter) {
				mc.EXPECT().GetUsersContext(gomock.Any()).Return(nil, errors.New("i don't think so"))
			},
			nil,
			true,
		},
		{
			"zero users",
			fields{options: DefOptions},
			args{context.Background()},
			func(mc *mockClienter) {
				mc.EXPECT().GetUsersContext(gomock.Any()).Return([]slack.User{}, nil)
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := NewmockClienter(gomock.NewController(t))

			tt.expectFn(mc)

			sd := &Session{
				client: mc,
				cfg:    tt.fields.options,
			}
			got, err := sd.fetchUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.fetchUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Session.fetchUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetUsers(t *testing.T) {
	dir := t.TempDir()
	type fields struct {
		options Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expectFn func(*mockClienter)
		want     types.Users
		wantErr  bool
	}{
		{
			"everything goes as planned",
			fields{options: Config{
				UserCache: CacheConfig{Filename: gimmeTempFile(t, dir), Retention: 5 * time.Hour},
				Limits: Limits{
					Tier2: TierLimits{Burst: 1},
					Tier3: TierLimits{Burst: 1},
				},
			}},
			args{context.Background()},
			func(mc *mockClienter) {
				mc.EXPECT().GetUsersContext(gomock.Any()).Return([]slack.User(testUsers), nil)
			},
			testUsers,
			false,
		},
		{
			"loaded from cache",
			fields{options: Config{
				UserCache: CacheConfig{Filename: gimmeTempFileWithUsers(t, dir), Retention: 5 * time.Hour},
				Limits: Limits{
					Tier2: TierLimits{Burst: 1},
					Tier3: TierLimits{Burst: 1},
				},
			}},
			args{context.Background()},
			func(mc *mockClienter) {
				// we don't expect any API calls
			},
			testUsers,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := NewmockClienter(gomock.NewController(t))

			tt.expectFn(mc)

			sd := &Session{
				client:  mc,
				wspInfo: &slack.AuthTestResponse{TeamID: testSuffix},
				cfg:     tt.fields.options,
				log:     logger.Silent,
			}
			got, err := sd.GetUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Session.GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func gimmeTempFile(t *testing.T, dir string) string {
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Errorf("error closing test file: %s", err)
	}
	return f.Name()
}

func gimmeTempFileWithUsers(t *testing.T, dir string) string {
	f := gimmeTempFile(t, dir)
	m, err := cache.NewManager("", cache.WithUserCacheBase(f))
	if err != nil {
		t.Fatal(err)
	}
	if err := m.SaveUsers(testSuffix, testUsers); err != nil {
		t.Fatal(err)
	}
	return f
}
