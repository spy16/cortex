package chunk

import (
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

type TodoData struct {
	Deadline time.Time  `json:"deadline"`
	Items    []TodoItem `json:"items"`
}

type TodoItem struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (data TodoData) Kind() string { return KindTodo }

func (data *TodoData) Validate() error {
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
