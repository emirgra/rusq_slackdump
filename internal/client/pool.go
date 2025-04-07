package client

import (
	"context"
	"io"
	"sync"

	"github.com/rusq/slack"
)

type Pool struct {
	pool []SlackClienter
	mu   sync.Mutex
	strategy
}

// NewPool wraps the slack.Client with the edge client, so that the edge
// client can be used as a fallback.
func NewPool(scl ...SlackClienter) *Pool {
	return &Pool{
		pool:     scl,
		strategy: newRoundRobin(len(scl)),
	}
}

// next returns the next client in the pool using the current strategy.
func (p *Pool) next() SlackClienter {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.pool) == 0 {
		panic("no clients in pool")
	}
	return p.pool[p.strategy.next()]
}

func (p *Pool) Client() (*slack.Client, bool) {
	return p.next().Client()
}

func (p *Pool) AuthTestContext(ctx context.Context) (response *slack.AuthTestResponse, err error) {
	return p.next().AuthTestContext(ctx)
}

func (p *Pool) GetConversationHistoryContext(ctx context.Context, params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error) {
	return p.next().GetConversationHistoryContext(ctx, params)
}

func (p *Pool) GetConversationRepliesContext(ctx context.Context, params *slack.GetConversationRepliesParameters) (msgs []slack.Message, hasMore bool, nextCursor string, err error) {
	return p.next().GetConversationRepliesContext(ctx, params)
}

func (p *Pool) GetUsersPaginated(options ...slack.GetUsersOption) slack.UserPagination {
	return p.next().GetUsersPaginated(options...)
}

func (p *Pool) GetStarredContext(ctx context.Context, params slack.StarsParameters) ([]slack.StarredItem, *slack.Paging, error) {
	return p.next().GetStarredContext(ctx, params)
}

func (p *Pool) ListBookmarks(channelID string) ([]slack.Bookmark, error) {
	return p.next().ListBookmarks(channelID)
}

func (p *Pool) GetConversationsContext(ctx context.Context, params *slack.GetConversationsParameters) (channels []slack.Channel, nextCursor string, err error) {
	return p.next().GetConversationsContext(ctx, params)
}

func (p *Pool) GetConversationInfoContext(ctx context.Context, input *slack.GetConversationInfoInput) (*slack.Channel, error) {
	return p.next().GetConversationInfoContext(ctx, input)
}

func (p *Pool) GetUsersInConversationContext(ctx context.Context, params *slack.GetUsersInConversationParameters) ([]string, string, error) {
	return p.next().GetUsersInConversationContext(ctx, params)
}

func (p *Pool) GetFileContext(ctx context.Context, downloadURL string, writer io.Writer) error {
	return p.next().GetFileContext(ctx, downloadURL, writer)
}

func (p *Pool) GetUsersContext(ctx context.Context, options ...slack.GetUsersOption) ([]slack.User, error) {
	return p.next().GetUsersContext(ctx, options...)
}

func (p *Pool) GetEmojiContext(ctx context.Context) (map[string]string, error) {
	return p.next().GetEmojiContext(ctx)
}

func (p *Pool) SearchMessagesContext(ctx context.Context, query string, params slack.SearchParameters) (*slack.SearchMessages, error) {
	return p.next().SearchMessagesContext(ctx, query, params)
}

func (p *Pool) SearchFilesContext(ctx context.Context, query string, params slack.SearchParameters) (*slack.SearchFiles, error) {
	return p.next().SearchFilesContext(ctx, query, params)
}

func (p *Pool) GetFileInfoContext(ctx context.Context, fileID string, count int, page int) (*slack.File, []slack.Comment, *slack.Paging, error) {
	return p.next().GetFileInfoContext(ctx, fileID, count, page)
}

func (p *Pool) GetUserInfoContext(ctx context.Context, user string) (*slack.User, error) {
	return p.next().GetUserInfoContext(ctx, user)
}
