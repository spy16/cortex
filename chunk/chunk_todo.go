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

func (data TodoData) Kind() string { return KindTodo }

func (data *TodoData) Validate() error {
	return nil
}
