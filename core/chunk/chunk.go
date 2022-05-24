package chunk

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

// Supported types of chunks.
const (
	KindNote  = "NOTE"
	KindTodo  = "TODO"
	KindImage = "IMAGE"
)

const (
	genIDLen   = 6
	genIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Chunk represents a piece of information written down by a user.
type Chunk struct {
	ID        string    `json:"id"`
	Data      Data      `json:"data"`
	Tags      []string  `json:"tags"`
	Author    string    `json:"author"`
	Parent    string    `json:"parent"`
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
	Tags   []string `json:"tags"`
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

	c.Author = strings.TrimSpace(c.Author)
	if c.Author == "" {
		return errors.ErrInvalid.WithMsgf("author must be set")
	}

	return c.genID()
}

func (c *Chunk) Apply(upd Updates) {
	if upd.Parent != "" {
		c.Parent = upd.Parent
	}

	if upd.Tags != nil {
		c.Tags = cleanTags(upd.Tags)
	}
}

func (c *Chunk) genID() error {
	kind := c.Data.Kind()
	idPrefix := strings.ToLower(string(kind[0]))
	idSuffix := randString(genIDLen, genIDChars)

	c.ID = fmt.Sprintf("%s-%s", idPrefix, idSuffix)
	return nil
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

func randString(n int, charset string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
