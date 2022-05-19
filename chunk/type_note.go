package chunk

type NoteData struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

func (n NoteData) Type() string { return TypeNote }

func (n *NoteData) Validate() error {
	return nil
}
