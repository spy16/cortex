package chunk

import (
	"net/url"
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

type ImageData struct {
	URL     string `json:"url"`
	Alt     string `json:"alt,omitempty"`
	Caption string `json:"caption,omitempty"`
}

func (data ImageData) Kind() string { return KindImage }

func (data *ImageData) Validate() error {
	data.URL = strings.TrimSpace(data.URL)
	data.Alt = strings.TrimSpace(data.Alt)
	data.Caption = strings.TrimSpace(data.Caption)

	if data.URL == "" {
		return errors.ErrInvalid.WithMsgf("image url must be specified")
	} else if _, err := url.Parse(data.URL); err != nil {
		return errors.ErrInvalid.WithMsgf("invalid image url").WithCausef(err.Error())
	}
	return nil
}
