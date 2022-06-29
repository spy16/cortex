package kind

import (
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

type TodoData struct {
	Items    []TodoItem `json:"items"`
	Deadline int64      `json:"deadline,omitempty"`
}

type TodoItem struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (data *TodoData) ValidateData() error {
	items := data.Items
	data.Items = nil
	for _, item := range items {
		item.Text = strings.TrimSpace(item.Text)
		if item.Text != "" {
			data.Items = append(data.Items, item)
		}
	}

	if len(data.Items) == 0 {
		return errors.ErrInvalid.WithMsgf("todo list must have at-least 1 item")
	}

	return nil
}
