package kind

import (
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

const (
	FormatMarkdown = "markdown"
)

type NoteData struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

func (data *NoteData) ValidateData() error {
	data.Text = strings.TrimSpace(data.Text)
	data.Format = strings.ToLower(strings.TrimSpace(data.Format))

	if data.Format == "" || data.Format == FormatMarkdown {
		data.Format = FormatMarkdown
	}

	if data.Text == "" {
		return errors.ErrInvalid.WithMsgf("data must contain text field")
	}
	return nil
}
