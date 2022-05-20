package chunk

type ImageData struct {
	Alt     string `json:"alt"`
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

func (data ImageData) Kind() string { return KindImage }

func (data *ImageData) Validate() error {
	return nil
}
