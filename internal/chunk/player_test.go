package chunk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/rusq/slackdump/v2/internal/chunk/state"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

var testThreads = []Chunk{
	{
		Type:      CThreadMessages,
		Timestamp: 1234567890,
		ChannelID: "C1234567890",
		IsThread:  true,
		Count:     2,
		Parent: &slack.Message{
			Msg: slack.Msg{
				ThreadTimestamp: "1234567890.123456",
			},
		},
		Messages: []slack.Message{
			{
				Msg: slack.Msg{
					ThreadTimestamp: "1234567890.123456",
					Timestamp:       "1234567890.123456",
					Text:            "Hello, world!",
				},
			},
			{
				Msg: slack.Msg{
					ThreadTimestamp: "1234567890.123456",
					Timestamp:       "1234567890.123457",
					Text:            "Hello, Slack!",
				},
			},
		},
	},
	{
		Type:      CThreadMessages,
		Timestamp: 1234567891,
		ChannelID: "C1234567890",
		IsThread:  true,
		Count:     2,
		Parent: &slack.Message{
			Msg: slack.Msg{
				ThreadTimestamp: "1234567890.123458",
			},
		},
		Messages: []slack.Message{
			{
				Msg: slack.Msg{
					ThreadTimestamp: "1234567890.123458",
					Timestamp:       "1234567890.200000",
					Text:            "Hello, world!",
				},
			},
			{
				Msg: slack.Msg{
					ThreadTimestamp: "1234567890.123458",
					Timestamp:       "1234567890.300000",
					Text:            "Hello, Slack!",
				},
			},
		},
	},
	{
		Type:      CThreadMessages,
		Timestamp: 1234567890,
		ChannelID: "C1234567890",
		IsThread:  true,
		Count:     2,
		Parent: &slack.Message{
			Msg: slack.Msg{
				ThreadTimestamp: "1234567890.123456",
			},
		},
		Messages: []slack.Message{
			{
				Msg: slack.Msg{
					ThreadTimestamp: "1234567890.123456",
					Timestamp:       "1234567890.400000",
					Text:            "Hello again world",
				},
			},
			{
				Msg: slack.Msg{
					ThreadTimestamp: "1234567890.123456",
					Timestamp:       "1234567890.500000",
					Text:            "Hello again Slack!",
				},
			},
		},
	},
}

var testThreadsIndex = index{
	"tC1234567890:1234567890.123456": []int64{0, 1225},
	"tC1234567890:1234567890.123458": []int64{612},
}

var testChunks = []Chunk{
	{Type: CChannelInfo, ChannelID: "C1234567890", Channel: &slack.Channel{GroupConversation: slack.GroupConversation{Conversation: slack.Conversation{ID: "C1234567890"}}}},
	{Type: CMessages, ChannelID: "C1234567890", Messages: []slack.Message{
		{Msg: slack.Msg{Timestamp: "1234567890.100000", Text: "message1"}},
		{Msg: slack.Msg{Timestamp: "1234567890.200000", Text: "message2"}},
		{Msg: slack.Msg{Timestamp: "1234567890.300000", Text: "message3"}},
		{Msg: slack.Msg{Timestamp: "1234567890.400000", Text: "message4"}},
		{Msg: slack.Msg{Timestamp: "1234567890.500000", Text: "message5"}},
	}},
	{Type: CMessages, ChannelID: "C1234567890", Messages: []slack.Message{
		{Msg: slack.Msg{Timestamp: "1234567890.600000", Text: "Hello, again!"}},
		{Msg: slack.Msg{Timestamp: "1234567890.700000", Text: "And again!"}},
	}},
	{Type: CMessages, ChannelID: "C1234567890", Messages: []slack.Message{
		{Msg: slack.Msg{Timestamp: "1234567890.800000", Text: "And again!"}},
		{
			Msg: slack.Msg{
				ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa70",
				ThreadTimestamp: "1234567890.800000",
				Timestamp:       "1234567890.800000",
				Text:            "parent message",
			},
		},
	}},
	{
		Type:      CThreadMessages,
		ChannelID: "C1234567890",
		IsThread:  true,
		Parent: &slack.Message{
			Msg: slack.Msg{
				ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa70",
				ThreadTimestamp: "1234567890.800000",
				Timestamp:       "1234567890.800000",
				Text:            "parent message",
			},
		},
		Timestamp: 1234567890,
		Count:     2,
		Messages: []slack.Message{
			{
				Msg: slack.Msg{
					ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa71",
					Timestamp:       "1234567890.900000",
					ThreadTimestamp: "1234567890.900000",
					Text:            "Hello, world!",
				},
			},
			{
				Msg: slack.Msg{
					ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa72",
					Timestamp:       "1234567891.100000",
					ThreadTimestamp: "1234567890.123456",
					Text:            "Hello, Slack!",
				},
			},
		},
	},
	// chunks from another channel
	{Type: CChannelInfo, ChannelID: "C987654321", Channel: &slack.Channel{GroupConversation: slack.GroupConversation{Conversation: slack.Conversation{ID: "C987654321"}}}},
	{Type: CMessages, ChannelID: "C987654321", Messages: []slack.Message{
		{Msg: slack.Msg{Timestamp: "1234567890.100000", Text: "message1"}},
		{Msg: slack.Msg{Timestamp: "1234567890.200000", Text: "message2"}},
		{Msg: slack.Msg{Timestamp: "1234567890.300000", Text: "message3"}},
		{Msg: slack.Msg{Timestamp: "1234567890.400000", Text: "message4"}},
		{Msg: slack.Msg{Timestamp: "1234567890.500000", Text: "message5"}},
	}},
	{Type: CMessages, ChannelID: "C987654321", Messages: []slack.Message{
		{Msg: slack.Msg{Timestamp: "1234567890.600000", Text: "Hello, again!"}},
		{Msg: slack.Msg{Timestamp: "1234567890.700000", Text: "And again!"}},
	}},
	{Type: CMessages, ChannelID: "C987654321", Messages: []slack.Message{
		{Msg: slack.Msg{Timestamp: "1234567890.800000", Text: "And again!"}},
		{
			Msg: slack.Msg{
				ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa70",
				ThreadTimestamp: "1234567890.800000",
				Timestamp:       "1234567890.800000",
				Text:            "parent message",
			},
		},
	}},
	{
		Type:      CThreadMessages,
		ChannelID: "C987654321",
		IsThread:  true,
		Parent: &slack.Message{
			Msg: slack.Msg{
				ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa70",
				ThreadTimestamp: "1234567890.800000",
				Timestamp:       "1234567890.800000",
				Text:            "parent message",
			},
		},
		Timestamp: 1234567890,
		Count:     2,
		Messages: []slack.Message{
			{
				Msg: slack.Msg{
					ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa71",
					Timestamp:       "1234567890.900000",
					ThreadTimestamp: "1234567890.900000",
					Text:            "Hello, world!",
				},
			},
			{
				Msg: slack.Msg{
					ClientMsgID:     "ec821bf2-c241-471d-b511-967b6ed4aa72",
					Timestamp:       "1234567891.100000",
					ThreadTimestamp: "1234567890.123456",
					Text:            "Hello, Slack!",
				},
			},
		},
	},
}

func marshalEvents(t *testing.T, v []Chunk) []byte {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, e := range v {
		if err := enc.Encode(e); err != nil {
			t.Fatal(err)
		}
	}
	return buf.Bytes()
}

func Test_indexRecords(t *testing.T) {
	type args struct {
		rs io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    index
		wantErr bool
	}{
		{
			name: "single thread",
			args: args{
				rs: bytes.NewReader(marshalEvents(t, testThreads)),
			},
			want:    testThreadsIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := indexRecords(json.NewDecoder(tt.args.rs))
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("indexRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_Thread(t *testing.T) {
	data := marshalEvents(t, testThreads)
	p := Player{
		rs:      bytes.NewReader(data),
		idx:     testThreadsIndex,
		pointer: make(offsets),
	}
	m, err := p.Thread("C1234567890", "1234567890.123456")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(m))
	}
	// again
	m, err = p.Thread("C1234567890", "1234567890.123456")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(m))
	}
	// should error
	m, err = p.Thread("C1234567890", "1234567890.123456")
	if !errors.Is(err, io.EOF) {
		t.Error(err, "expected io.EOF")
	}
	if len(m) > 0 {
		t.Fatalf("expected 0 messages, got %d", len(m))
	}
}

func TestPlayer_FileState(t *testing.T) {
	type fields struct {
		rs         io.ReadSeeker
		pointer    offsets
		idx        index
		lastOffset atomic.Int64
	}
	tests := []struct {
		name    string
		fields  fields
		want    *state.State
		wantErr bool
	}{
		{
			name: "single thread",
			fields: fields{
				rs: bytes.NewReader(marshalEvents(t, testThreads)),
			},
			want: &state.State{
				Version:  state.Version,
				Channels: make(map[string]int64),
				Threads: map[string]int64{
					"C1234567890:1234567890.123456": 1234567890500000,
					"C1234567890:1234567890.123458": 1234567890300000,
				},
				Files: make(map[string]string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				rs:         tt.fields.rs,
				pointer:    tt.fields.pointer,
				idx:        tt.fields.idx,
				lastOffset: tt.fields.lastOffset,
			}
			got, err := p.State()
			if (err != nil) != tt.wantErr {
				t.Errorf("Player.FileState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Player.FileState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_AllChannels(t *testing.T) {
	type fields struct {
		rs         io.ReadSeeker
		pointer    offsets
		lastOffset atomic.Int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "ok",
			fields: fields{
				rs: bytes.NewReader(marshalEvents(t, testChunks)),
			},
			want: []string{"C1234567890", "C987654321"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, err := indexRecords(json.NewDecoder(tt.fields.rs))
			if err != nil {
				t.Fatal(err)
			}
			p := &Player{
				rs:         tt.fields.rs,
				idx:        idx,
				pointer:    tt.fields.pointer,
				lastOffset: tt.fields.lastOffset,
			}
			if got := p.AllChannels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.AllChannels() = %v, want %v", got, tt.want)
			}
		})
	}
}
