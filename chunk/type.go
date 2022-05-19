package chunk

import (
	"encoding/json"
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

// Supported types of chunks.
const (
	TypeNote  = "NOTE"
	TypeTodo  = "TODO"
	TypeImage = "IMAGE"
)

// Data implementations represent the data of different types of chunks.
type Data interface {
	Type() string
	Validate() error
}

// ParseData parses the given raw data of specified chunkType into the Data
// implementation.
func ParseData(chunkType, data string) (Data, error) {
	chunkType = strings.ToUpper(chunkType)

	var into Data
	switch chunkType {
	case TypeNote:
		into = &NoteData{}

	case TypeTodo:
		into = &TodoData{}

	case TypeImage:
		into = &ImageData{}

	default:
		return nil, errors.ErrInvalid.WithMsgf("invalid chunk-type '%s'", chunkType)
	}

	if err := json.Unmarshal([]byte(data), into); err != nil {
		return nil, errors.ErrInvalid.
			WithMsgf("failed to interpret data as '%s'", chunkType).
			WithCausef(err.Error())
	}
	return into, nil
}
