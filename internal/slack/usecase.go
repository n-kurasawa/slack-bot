package slack

import (
	"fmt"

	slackapi "github.com/slack-go/slack"
)

type UseCase struct {
	client   SlackClient
	imgStore ImageStore
}

type Image struct {
	ID   int
	URL  string
	Name string
}

func NewUseCase(client SlackClient, store ImageStore) *UseCase {
	return &UseCase{
		client:   client,
		imgStore: store,
	}
}

func (u *UseCase) SendHelloWorld(channelID string) error {
	_, _, err := u.client.PostMessage(
		channelID,
		slackapi.MsgOptionText("world", false),
	)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗: %w", err)
	}
	return nil
}

func (u *UseCase) SendImage(channelID string, name string) error {
	var img *Image
	var err error

	if name != "" {
		// 名前指定がある場合
		img, err = u.imgStore.GetImageByName(name)
	} else {
		// 名前指定がない場合はランダム
		img, err = u.imgStore.GetImage()
	}

	if err != nil {
		return fmt.Errorf("画像の取得に失敗: %w", err)
	}

	_, _, err = u.client.PostMessage(
		channelID,
		slackapi.MsgOptionText(fmt.Sprintf("%s\n%s", img.Name, img.URL), false),
	)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗: %w", err)
	}
	return nil
}

func (u *UseCase) SaveImage(channelID, name, url string) error {
	if err := u.imgStore.SaveImage(name, url); err != nil {
		return fmt.Errorf("画像の保存に失敗: %w", err)
	}

	_, _, err := u.client.PostMessage(
		channelID,
		slackapi.MsgOptionText("画像を保存しました :white_check_mark:", false),
	)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗: %w", err)
	}
	return nil
}

func (u *UseCase) SendInvalidCommandError(channelID string) error {
	_, _, err := u.client.PostMessage(
		channelID,
		slackapi.MsgOptionText("不正なコマンド形式です。使用方法: updateImage NAME URL", false),
	)
	if err != nil {
		return fmt.Errorf("エラーメッセージの送信に失敗: %w", err)
	}
	return nil
}
