package bot

import "github.com/slack-go/slack"

type ImageStore interface {
	GetImage() (*Image, error)
	GetImageByName(name string) (*Image, error)
	SaveImage(name, url string) error
}

type Image struct {
	ID   int
	URL  string
	Name string
}

type SlackClient interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}
