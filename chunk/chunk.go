package chunk

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

const maxIDLen = 12

// Supported types of chunks.
const (
	KindNote  = "NOTE"
	KindTodo  = "TODO"
	KindImage = "IMAGE"
)

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// Chunk represents a piece of information written down by a user.
type Chunk struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Data      Data      `json:"data"`
	Tags      []string  `json:"tags"`
	Rank      string    `json:"rank"`
	Author    string    `json:"author"`
	Parent    string    `json:"parent,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Data implementations represent the data of different types of chunks.
type Data interface {
	Kind() string
	Validate() error
}

// Updates represents modifications to a chunk. Zero-values are treated
// as no-update.
type Updates struct {
	Data   Data     `json:"data"`
	Tags   []string `json:"tags"`
	Rank   string   `json:"rank"`
	Parent string   `json:"parent"`
}

// ParseData parses the given raw data of specified chunkType into the Data
// implementation.
func ParseData(kind, data string) (Data, error) {
	kind = strings.ToUpper(kind)

	var into Data
	switch kind {
	case KindNote:
		into = &NoteData{}

	case KindTodo:
		into = &TodoData{}

	case KindImage:
		into = &ImageData{}

	default:
		return nil, errors.ErrInvalid.WithMsgf("invalid kind '%s'", kind)
	}

	if err := json.Unmarshal([]byte(data), into); err != nil {
		return nil, errors.ErrInvalid.
			WithMsgf("failed to interpret data as '%s'", kind).
			WithCausef(err.Error())
	}
	return into, nil
}

func (c *Chunk) Validate() error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
		c.UpdatedAt = c.CreatedAt
	}

	if c.Data == nil {
		return errors.ErrInvalid.WithMsgf("chunk data must be set")
	} else if err := c.Data.Validate(); err != nil {
		return err
	}
	c.Kind = c.Data.Kind()

	c.ID = strings.TrimSpace(c.ID)
	if err := validateID(c.ID, false); err != nil {
		return err
	}

	c.Parent = strings.TrimSpace(c.Parent)
	if err := validateID(c.Parent, true); err != nil {
		return err
	}

	c.Author = strings.TrimSpace(c.Author)
	if c.Author == "" {
		return errors.ErrInvalid.WithMsgf("author must be set")
	}

	return nil
}

func validateID(id string, optional bool) error {
	if optional && id == "" {
		return nil
	}

	if len(id) == 0 || len(id) > maxIDLen {
		return errors.ErrInvalid.WithMsgf("id must have 1-%d characters", maxIDLen)
	} else if !idPattern.MatchString(id) {
		return errors.ErrInvalid.WithMsgf("id must match pattern '%s'", idPattern)
	}
	return nil
}
