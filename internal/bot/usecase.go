package bot

import (
	"fmt"

	"github.com/slack-go/slack"
)

type MessageEventService struct {
	client   *slack.Client
	imgStore ImageStore
}

func NewMessageEventService(client *slack.Client, store ImageStore) *MessageEventService {
	return &MessageEventService{
		client:   client,
		imgStore: store,
	}
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
