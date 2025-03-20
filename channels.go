package slackdump

// In this file: channel/conversations and thread related code.

import (
	"context"
	"runtime/trace"
	"time"

	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3/internal/network"
	"github.com/rusq/slackdump/v3/types"
)

// GetChannels list all conversations for a user.  `chanTypes` specifies the
// type of messages to fetch.  See github.com/rusq/slack docs for possible
// values.  If large number of channels is to be returned, consider using
// StreamChannels.
func (s *Session) GetChannels(ctx context.Context, chanTypes ...string) (types.Channels, error) {
	var allChannels types.Channels
	if err := s.getChannels(ctx, chanTypes, func(cc types.Channels) error {
		allChannels = append(allChannels, cc...)
		return nil
	}); err != nil {
		return allChannels, err
	}
	return allChannels, nil
}

// StreamChannels requests the channels from the API and calls the callback
// function cb for each.
func (s *Session) StreamChannels(ctx context.Context, chanTypes []string, cb func(ch slack.Channel) error) error {
	return s.getChannels(ctx, chanTypes, func(chans types.Channels) error {
		for _, ch := range chans {
			if err := cb(ch); err != nil {
				return err
			}
		}
		return nil
	})
}

// getChannels list all conversations for a user.  `chanTypes` specifies
// the type of messages to fetch.  See github.com/rusq/slack docs for possible
// values
func (s *Session) getChannels(ctx context.Context, chanTypes []string, cb func(types.Channels) error) error {
	ctx, task := trace.NewTask(ctx, "getChannels")
	defer task.End()

	limiter := s.limiter(network.Tier2)

	if chanTypes == nil {
		chanTypes = AllChanTypes
	}

	params := &slack.GetConversationsParameters{Types: chanTypes, Limit: s.cfg.limits.Request.Channels}
	fetchStart := time.Now()
	var total int
	for i := 1; ; i++ {
		var (
			chans   []slack.Channel
			nextcur string
		)
		reqStart := time.Now()
		if err := network.WithRetry(ctx, limiter, s.cfg.limits.Tier3.Retries, func(ctx context.Context) error {
			var err error
			trace.WithRegion(ctx, "GetConversationsContext", func() {
				chans, nextcur, err = s.client.GetConversationsContext(ctx, params)
			})
			return err
		}); err != nil {
			return err
		}

		if err := cb(chans); err != nil {
			return err
		}
		total += len(chans)

		s.log.InfoContext(ctx, "channels", "request", i, "fetched", len(chans), "total", total,
			"speed", float64(len(chans))/time.Since(reqStart).Seconds(),
			"avg", float64(total)/time.Since(fetchStart).Seconds(),
		)

		if nextcur == "" {
			s.log.InfoContext(ctx, "channels fetch complete", "total", total)
			break
		}

		params.Cursor = nextcur

		if err := limiter.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

// GetChannelMembers returns a list of all members in a channel.
func (sd *Session) GetChannelMembers(ctx context.Context, channelID string) ([]string, error) {
	var ids []string
	var cursor string
	for {
		var uu []string
		var next string
		if err := network.WithRetry(ctx, sd.limiter(network.Tier4), sd.cfg.limits.Tier4.Retries, func(ctx context.Context) error {
			var err error
			uu, next, err = sd.client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
				ChannelID: channelID,
				Cursor:    cursor,
			})
			return err
		}); err != nil {
			return nil, err
		}
		ids = append(ids, uu...)

		if next == "" {
			break
		}
		cursor = next
	}
	return ids, nil
}
