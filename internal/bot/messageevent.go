package bot

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type MessageEventService struct {
	client   SlackClient
	imgStore ImageStore
}

func NewMessageEventService(client SlackClient, store ImageStore) *MessageEventService {
	return &MessageEventService{
		client:   client,
		imgStore: store,
	}
}

func (s *MessageEventService) HandleMessage(event *slackevents.MessageEvent) error {
	switch {
	case event.Text == "hello":
		return s.SendHelloWorld(event.Channel)

	case strings.HasPrefix(event.Text, "image"):
		parts := strings.Fields(event.Text)
		var name string
		if len(parts) > 1 {
			name = parts[1]
		}
		return s.SendImage(event.Channel, name)

	case strings.HasPrefix(event.Text, "updateImage "):
		parts := strings.Fields(event.Text)
		if len(parts) != 3 {
			return s.SendInvalidCommandError(event.Channel)
		}

		name := parts[1]
		url := parts[2]
		return s.SaveImage(event.Channel, name, url)
	}

	return nil
}

func (s *MessageEventService) SendHelloWorld(channelID string) error {
	_, _, err := s.client.PostMessage(
		channelID,
		slack.MsgOptionText("world", false),
	)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗: %w", err)
	}
	return nil
}

func (s *MessageEventService) SendImage(channelID string, name string) error {
	var img *Image
	var err error

	if name != "" {
		// 名前指定がある場合
		img, err = s.imgStore.GetImageByName(name)
	} else {
		// 名前指定がない場合はランダム
		img, err = s.imgStore.GetImage()
	}

	if err != nil {
		return fmt.Errorf("画像の取得に失敗: %w", err)
	}

	_, _, err = s.client.PostMessage(
		channelID,
		slack.MsgOptionText(fmt.Sprintf("%s\n%s", img.Name, img.URL), false),
	)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗: %w", err)
	}
	return nil
}

func (s *MessageEventService) SaveImage(channelID, name, url string) error {
	if err := s.imgStore.SaveImage(name, url); err != nil {
		return fmt.Errorf("画像の保存に失敗: %w", err)
	}

	_, _, err := s.client.PostMessage(
		channelID,
		slack.MsgOptionText("画像を保存しました :white_check_mark:", false),
	)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗: %w", err)
	}
	return nil
}

func (s *MessageEventService) SendInvalidCommandError(channelID string) error {
	_, _, err := s.client.PostMessage(
		channelID,
		slack.MsgOptionText("不正なコマンド形式です。使用方法: updateImage NAME URL", false),
	)
	if err != nil {
		return fmt.Errorf("エラーメッセージの送信に失敗: %w", err)
	}
	return nil
}
