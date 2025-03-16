package bot

import "github.com/slack-go/slack"

// SlackClient はSlack APIとの通信を抽象化するインターフェースです
type SlackClient interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}
