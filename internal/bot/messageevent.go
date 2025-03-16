package bot

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack/slackevents"
)

type MessageEventService struct {
	imgStore ImageStore
}

type Response struct {
	Text string
}

func NewMessageEventService(store ImageStore) *MessageEventService {
	return &MessageEventService{
		imgStore: store,
	}
}

func (s *MessageEventService) HandleMessage(event *slackevents.MessageEvent) (*Response, error) {
	switch {
	case event.Text == "hello":
		return &Response{
			Text: "world",
		}, nil

	case strings.HasPrefix(event.Text, "image"):
		parts := strings.Fields(event.Text)
		var name string
		if len(parts) > 1 {
			name = parts[1]
		}

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
			return nil, fmt.Errorf("画像の取得に失敗: %w", err)
		}

		return &Response{
			Text: fmt.Sprintf("%s\n%s", img.Name, img.URL),
		}, nil

	case strings.HasPrefix(event.Text, "updateImage "):
		parts := strings.Fields(event.Text)
		if len(parts) != 3 {
			return &Response{
				Text: "不正なコマンド形式です。使用方法: updateImage NAME URL",
			}, nil
		}

		name := parts[1]
		url := parts[2]

		if err := s.imgStore.SaveImage(name, url); err != nil {
			return nil, fmt.Errorf("画像の保存に失敗: %w", err)
		}

		return &Response{
			Text: "画像を保存しました :white_check_mark:",
		}, nil
	}

	return nil, nil
}
