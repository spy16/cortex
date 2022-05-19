package chunk

type ImageData struct {
	Alt     string `json:"alt"`
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

func (i ImageData) Type() string { return TypeImage }

func (i ImageData) Validate() error {
	return nil
}
