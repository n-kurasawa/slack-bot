package bot

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack/slackevents"
)

type MessageEventHandler struct {
	imgStore ImageStore
}

func NewMessageEventHandler(store ImageStore) *MessageEventHandler {
	return &MessageEventHandler{
		imgStore: store,
	}
}

func (s *MessageEventHandler) HandleMessage(event *slackevents.MessageEvent) (string, error) {
	switch {
	case event.Text == "hello":
		return s.handleHello()

	case event.Text == "imageList":
		return s.handleImageList()

	case strings.HasPrefix(event.Text, "image"):
		return s.handleImage(event.Text)

	case strings.HasPrefix(event.Text, "updateImage "):
		return s.handleUpdateImage(event.Text)

	default:
		return "", nil
	}
}

// helloコマンドの処理
func (s *MessageEventHandler) handleHello() (string, error) {
	return "world", nil
}

// imageコマンドの処理
func (s *MessageEventHandler) handleImage(text string) (string, error) {
	parts := strings.Fields(text)
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
		return "", fmt.Errorf("画像の取得に失敗: %w", err)
	}

	return fmt.Sprintf("%s\n%s", img.Name, img.URL), nil
}

// updateImageコマンドの処理
func (s *MessageEventHandler) handleUpdateImage(text string) (string, error) {
	parts := strings.Fields(text)
	if len(parts) != 3 {
		return "不正なコマンド形式です。使用方法: updateImage NAME URL", nil
	}

	name := parts[1]
	url := parts[2]

	if err := s.imgStore.SaveImage(name, url); err != nil {
		return "", fmt.Errorf("画像の保存に失敗: %w", err)
	}

	return "画像を保存しました :white_check_mark:", nil
}

// 登録されている画像の一覧を取得して表示する
func (s *MessageEventHandler) handleImageList() (string, error) {
	images, err := s.imgStore.ListImages()
	if err != nil {
		return "", fmt.Errorf("画像一覧の取得に失敗: %w", err)
	}

	if len(images) == 0 {
		return "登録されている画像はありません", nil
	}

	var result strings.Builder
	result.WriteString("登録されている画像一覧:\n")

	for i, img := range images {
		result.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, img.Name, img.URL))
	}

	return result.String(), nil
}
