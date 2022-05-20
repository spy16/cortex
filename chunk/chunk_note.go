package chunk

type NoteData struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

func (data NoteData) Kind() string { return KindNote }

func (data *NoteData) Validate() error {
	return nil
}
