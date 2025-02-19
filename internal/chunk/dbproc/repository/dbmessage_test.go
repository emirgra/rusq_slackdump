package repository

import (
	"context"
	"encoding/json"
	"iter"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rusq/slack"
	"github.com/stretchr/testify/assert"

	"github.com/rusq/slackdump/v3/internal/chunk"
	"github.com/rusq/slackdump/v3/internal/fixtures"
)

func minifyJSON[T any](t *testing.T, s string) []byte {
	t.Helper()
	var a T
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		t.Fatalf("minifyJSON: %v", err)
	}
	b, err := marshal(a)
	if err != nil {
		t.Fatalf("minifyJSON: %v", err)
	}
	return b
}

func TestNewDBMessage(t *testing.T) {
	type args struct {
		dbchunkID int64
		idx       int
		channelID string
		msg       *slack.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *DBMessage
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				dbchunkID: 100,
				idx:       222,
				channelID: "C123",
				msg:       fixtures.Load[*slack.Message](fixtures.SimpleMessageJSON),
			},
			want: &DBMessage{
				ID:        1645095505023899,
				ChunkID:   100,
				ChannelID: "C123",
				TS:        "1645095505.023899",
				IsParent:  false,
				Index:     222,
				NumFiles:  0,
				Text:      "Test message with Html chars &lt; &gt;",
				Data:      minifyJSON[slack.Message](t, fixtures.SimpleMessageJSON),
			},
			wantErr: false,
		},
		{
			name: "bot thread parent message",
			args: args{
				dbchunkID: 100,
				idx:       222,
				channelID: "C123",
				msg:       fixtures.Load[*slack.Message](fixtures.BotMessageThreadParentJSON),
			},
			want: &DBMessage{
				ID:        1648085300726649,
				ChunkID:   100,
				ChannelID: "C123",
				TS:        "1648085300.726649",
				ParentID:  ptr[int64](1648085300726649),
				ThreadTS:  ptr("1648085300.726649"),
				IsParent:  true,
				Index:     222,
				NumFiles:  0,
				Text:      "This content can't be displayed.",
				Data:      minifyJSON[slack.Message](t, fixtures.BotMessageThreadParentJSON),
			},
		},
		{
			name: "bot thread child message w files",
			args: args{
				dbchunkID: 100,
				idx:       222,
				channelID: "C123",
				msg:       fixtures.Load[*slack.Message](fixtures.BotMessageThreadChildJSON),
			},
			want: &DBMessage{
				ID:        1648085301269949,
				ChunkID:   100,
				ChannelID: "C123",
				TS:        "1648085301.269949",
				ParentID:  ptr[int64](1648085300726649),
				ThreadTS:  ptr("1648085300.726649"),
				IsParent:  false,
				Index:     222,
				NumFiles:  1,
				Text:      "",
				Data:      minifyJSON[slack.Message](t, fixtures.BotMessageThreadChildJSON),
			},
		},
		{
			name: "app message",
			args: args{
				dbchunkID: 100,
				idx:       222,
				channelID: "C123",
				msg:       fixtures.Load[*slack.Message](fixtures.AppMessageJSON),
			},
			want: &DBMessage{
				ID:        1586042786000100,
				ChunkID:   100,
				ChannelID: "C123",
				TS:        "1586042786.000100",
				IsParent:  false,
				Index:     222,
				NumFiles:  0,
				Text:      "",
				Data:      minifyJSON[slack.Message](t, fixtures.AppMessageJSON),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDBMessage(tt.args.dbchunkID, tt.args.idx, tt.args.channelID, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDBMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_messageRepository_Insert(t *testing.T) {
	// fixtures
	simpleDBMessage, err := NewDBMessage(1, 0, "C123", fixtures.Load[*slack.Message](fixtures.SimpleMessageJSON))
	if err != nil {
		t.Fatalf("newdbmessage: %v", err)
	}

	type args struct {
		ctx  context.Context
		conn PrepareExtContext
		m    *DBMessage
	}
	tests := []struct {
		name    string
		m       messageRepository
		args    args
		prepFn  utilityFn
		wantErr bool
		checkFn utilityFn
	}{
		{
			name: "ok",
			m:    messageRepository{},
			args: args{
				ctx:  context.Background(),
				conn: testConn(t),
				m:    simpleDBMessage,
			},
			prepFn:  prepChunk(chunk.CMessages),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepFn != nil {
				tt.prepFn(t, tt.args.conn)
			}
			m := NewMessageRepository()
			if err := m.Insert(tt.args.ctx, tt.args.conn, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFn != nil {
				tt.checkFn(t, tt.args.conn)
			}
		})
	}
}

func Test_messageRepository_InsertAll(t *testing.T) {
	type args struct {
		ctx   context.Context
		pconn PrepareExtContext
		mm    iter.Seq2[*DBMessage, error]
	}
	tests := []struct {
		name    string
		args    args
		prepFn  utilityFn
		want    int
		wantErr bool
		checkFn utilityFn
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				pconn: testConn(t),
				mm: toIter([]testResult[*DBMessage]{
					{V: &DBMessage{ID: 1, ChunkID: 1, ChannelID: "C123", TS: "1.1", IsParent: false, Index: 0, NumFiles: 0, Text: "test", Data: []byte(`{"text":"test"}`)}},
					toTestResult(NewDBMessage(1, 1, "C123", fixtures.Load[*slack.Message](fixtures.SimpleMessageJSON))),
				}),
			},
			prepFn:  prepChunk(chunk.CMessages),
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepFn != nil {
				tt.prepFn(t, tt.args.pconn)
			}
			m := NewMessageRepository()
			got, err := m.InsertAll(tt.args.ctx, tt.args.pconn, tt.args.mm)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.InsertAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("messageRepository.InsertAll() = %v, want %v", got, tt.want)
			}
			if tt.checkFn != nil {
				tt.checkFn(t, tt.args.pconn)
			}
		})
	}
}

var (
	msgA  = slack.Message{Msg: slack.Msg{Timestamp: "123.456", Text: "A"}}
	msgB  = slack.Message{Msg: slack.Msg{Timestamp: "124.555", Text: "B"}}
	msgB_ = slack.Message{Msg: slack.Msg{Timestamp: "124.555", Text: "B'"}}
	msgC  = slack.Message{Msg: slack.Msg{Timestamp: "125.777", Text: "C"}}

	msgX = slack.Message{Msg: slack.Msg{Timestamp: "123.456", Text: "X"}} // these belong to a different channel
	msgY = slack.Message{Msg: slack.Msg{Timestamp: "124.555", Text: "Y"}}
	msgZ = slack.Message{Msg: slack.Msg{Timestamp: "125.777", Text: "Z"}}

	dbmA  = must(NewDBMessage(1, 0, "C123", &msgA))
	dbmB  = must(NewDBMessage(1, 1, "C123", &msgB))
	dbmB_ = must(NewDBMessage(2, 0, "C123", &msgB_))
	dbmC  = must(NewDBMessage(2, 1, "C123", &msgC))

	dbmX = must(NewDBMessage(3, 0, "D124", &msgX))
	dbmY = must(NewDBMessage(3, 1, "D124", &msgY))
	dbmZ = must(NewDBMessage(4, 0, "D124", &msgZ))
)

func messagePrepFn(t *testing.T, conn PrepareExtContext) {
	// we will use 2 chunks, one old and one new for the same channel
	// they both will have 2 messages each, such as  (A, B),(B', C)
	// where B' will be an updated version of B.
	// Also, there are messages from a different channel, X, Y, Z.
	prepChunk(chunk.CMessages, chunk.CMessages, chunk.CMessages, chunk.CMessages)(t, conn)
	mr := NewMessageRepository()
	if err := mr.Insert(context.Background(), conn, dbmA, dbmB, dbmB_, dbmC, dbmX, dbmY, dbmZ); err != nil {
		t.Fatalf("insert: %v", err)
	}
}

func Test_messageRepository_Count(t *testing.T) {
	type fields struct {
		genericRepository genericRepository[DBMessage]
	}
	type args struct {
		ctx       context.Context
		conn      sqlx.QueryerContext
		channelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepFn  utilityFn
		want    int64
		wantErr bool
	}{
		{
			name: "count the most recent messages, without duplicates",
			fields: fields{
				genericRepository: genericRepository[DBMessage]{DBMessage{}},
			},
			args: args{
				ctx:       context.Background(),
				conn:      testConn(t),
				channelID: "C123",
			},
			prepFn:  messagePrepFn,
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepFn != nil {
				tt.prepFn(t, tt.args.conn.(PrepareExtContext))
			}
			r := messageRepository{
				genericRepository: tt.fields.genericRepository,
			}
			got, err := r.Count(tt.args.ctx, tt.args.conn, tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("messageRepository.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_messageRepository_AllForID(t *testing.T) {
	type fields struct {
		genericRepository genericRepository[DBMessage]
	}
	type args struct {
		ctx       context.Context
		conn      sqlx.QueryerContext
		channelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepFn  utilityFn
		want    []testResult[DBMessage]
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				genericRepository: genericRepository[DBMessage]{DBMessage{}},
			},
			args: args{
				ctx:       context.Background(),
				conn:      testConn(t),
				channelID: "C123",
			},
			prepFn: messagePrepFn,
			want: []testResult[DBMessage]{
				{V: *dbmA},
				{V: *dbmB_},
				{V: *dbmC},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepFn != nil {
				tt.prepFn(t, tt.args.conn.(PrepareExtContext))
			}
			r := messageRepository{
				genericRepository: tt.fields.genericRepository,
			}
			got, err := r.AllForID(tt.args.ctx, tt.args.conn, tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.AllForID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertIterResult(t, tt.want, got)
		})
	}
}

func Test_messageRepository_CountThread(t *testing.T) {
	type fields struct {
		genericRepository genericRepository[DBMessage]
	}
	type args struct {
		ctx       context.Context
		conn      sqlx.QueryerContext
		channelID string
		threadID  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepFn  utilityFn
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepFn != nil {
				tt.prepFn(t, tt.args.conn.(PrepareExtContext))
			}
			r := messageRepository{
				genericRepository: tt.fields.genericRepository,
			}
			got, err := r.CountThread(tt.args.ctx, tt.args.conn, tt.args.channelID, tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.CountThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("messageRepository.CountThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	tmAParent  = slack.Message{Msg: slack.Msg{Timestamp: "123.456", ThreadTimestamp: "123.456", Text: "A"}}
	tmBChannel = slack.Message{Msg: slack.Msg{Timestamp: "124.000", ThreadTimestamp: "123.456", Text: "B", SubType: "thread_broadcast"}}
	tmB        = slack.Message{Msg: slack.Msg{Timestamp: "124.000", ThreadTimestamp: "123.456", Text: "B", SubType: "thread_broadcast"}}
	tmC        = slack.Message{Msg: slack.Msg{Timestamp: "125.000", ThreadTimestamp: "123.456", Text: "C"}}
	tmD        = slack.Message{Msg: slack.Msg{Timestamp: "126.000", ThreadTimestamp: "123.456", Text: "D"}}
	tmC_       = slack.Message{Msg: slack.Msg{Timestamp: "125.000", ThreadTimestamp: "123.456", Text: "C'"}}

	dbtmAParent  = must(NewDBMessage(1, 0, "C123", &tmAParent))
	dbtmBChannel = must(NewDBMessage(1, 0, "C123", &tmBChannel))
	dbtmB        = must(NewDBMessage(2, 1, "C123", &tmB))
	dbtmC        = must(NewDBMessage(2, 1, "C123", &tmC))
	dbtmD        = must(NewDBMessage(2, 1, "C123", &tmD))
	dbtmC_       = must(NewDBMessage(3, 1, "C123", &tmC_))
)

func threadSetupFn(t *testing.T, conn PrepareExtContext) {
	// thread setup is the following:
	// chunk type_id subtype message   comment
	//     1       0    NULL       A   parent message
	//     1       0   bcast       B   thread broadcast in the channel - should not be included
	//     2       1   bcast       B   thread broadcast in the thread
	//     2       1    NULL       C   old thread message
	//     2       1    NULL       D   thread message
	//     3       1    NULL      C'   new thread message version of C.
	//
	//  The net result should be that we have 4 messages in the thread:
	//  A, B, C', D
	//
	//    chunk_id: 1                    2                    3
	prepChunk(chunk.CMessages, chunk.CThreadMessages, chunk.CThreadMessages)(t, conn)

	mr := NewMessageRepository()
	if err := mr.Insert(context.Background(), conn,
		dbtmAParent,
		dbtmBChannel,
		dbtmB,
		dbtmC,
		dbtmD,
		dbtmC_,
	); err != nil {
		t.Fatalf("insert: %v", err)
	}
}

func Test_messageRepository_AllForThread(t *testing.T) {
	type fields struct {
		genericRepository genericRepository[DBMessage]
	}
	type args struct {
		ctx       context.Context
		conn      sqlx.QueryerContext
		channelID string
		threadID  string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		preparFn utilityFn
		want     []testResult[DBMessage]
		wantErr  bool
	}{
		{
			name: "ok",
			fields: fields{
				genericRepository: genericRepository[DBMessage]{DBMessage{}},
			},
			args: args{
				ctx:       context.Background(),
				conn:      testConn(t),
				channelID: "C123",
				threadID:  "123.456",
			},
			preparFn: threadSetupFn,
			want: []testResult[DBMessage]{
				{V: *dbtmAParent},
				{V: *dbtmB},
				{V: *dbtmC_},
				{V: *dbtmD},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preparFn != nil {
				tt.preparFn(t, tt.args.conn.(PrepareExtContext))
			}
			r := messageRepository{
				genericRepository: tt.fields.genericRepository,
			}
			got, err := r.AllForThread(tt.args.ctx, tt.args.conn, tt.args.channelID, tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.AllForThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertIterResult(t, tt.want, got)
		})
	}
}
