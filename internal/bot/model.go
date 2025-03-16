package bot

type ImageStore interface {
	GetImage() (*Image, error)
	GetImageByName(name string) (*Image, error)
	SaveImage(name, url string) error
	ListImages() ([]Image, error)
}

type Image struct {
	ID   int
	URL  string
	Name string
}
