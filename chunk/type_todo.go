package chunk

import "time"

type TodoData struct {
	Deadline time.Time  `json:"deadline"`
	Items    []TodoItem `json:"items"`
}

type TodoItem struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (td *TodoData) Type() string { return TypeTodo }

func (td *TodoData) Validate() error {
	return nil
}
