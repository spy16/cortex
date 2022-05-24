package chunk

type UserData struct {
	Bio string `json:"bio,omitempty"`
}

func (u UserData) Kind() string { return KindUser }

func (u *UserData) Validate() error {
	return nil
}
