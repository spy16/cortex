package chunk

import (
	"context"
	"encoding/json"
	"math/rand"
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

const idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// Store implementation is responsible for persisting chunks.
type Store interface {
	Get(ctx context.Context, id string) (*Chunk, error)
	List(ctx context.Context, opts ListOptions) ([]Chunk, error)
	Create(ctx context.Context, c Chunk) error
	Update(ctx context.Context, id string, upd Updates) (*Chunk, error)
	Delete(ctx context.Context, id string) (*Chunk, error)
}

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
	} else if err := into.Validate(); err != nil {
		return nil, err
	}
	return into, nil
}

func (c *Chunk) Validate() error {
	c.Tags = cleanTags(c.Tags)
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

func (c *Chunk) Apply(upd Updates) {
	if upd.Data != nil {
		c.Kind = upd.Data.Kind()
		c.Data = upd.Data
	}

	if upd.Rank != "" {
		c.Rank = upd.Rank
	}

	if upd.Parent != "" {
		c.Parent = upd.Parent
	}

	if upd.Tags != nil {
		c.Tags = cleanTags(upd.Tags)
	}
}

func (c *Chunk) genID() { c.ID = randString(5, idCharset) }

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

func randString(n int, charset string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func cleanTags(tags []string) []string {
	set := map[string]struct{}{}

	var res []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}

		if _, exists := set[tag]; !exists {
			res = append(res, tag)
			set[tag] = struct{}{}
		}
	}
	return res
}
