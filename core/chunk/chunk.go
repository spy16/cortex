package chunk

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

const (
	genIDLen   = 6
	genIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Chunk represents a piece of information written down by a user.
type Chunk struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Data      string    `json:"data"`
	Tags      []string  `json:"tags"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KindRegistry is responsible for parsing, validating different kinds
// of chunk data.
type KindRegistry interface {
	Validate(kind, data string) error
}

// Data implementations represent the data of different types of chunks.
type Data interface {
	Kind() string
	Validate() error
}

// Updates represents modifications to a chunk. Zero-values are treated
// as no-update.
type Updates struct {
	Tags []string `json:"tags"`
	Kind string   `json:"kind"`
	Data string   `json:"data"`
}

func (c *Chunk) Validate() error {
	c.Tags = cleanTags(c.Tags)
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
		c.UpdatedAt = c.CreatedAt
	}

	if c.Kind == "" {
		return errors.ErrInvalid.WithMsgf("chunk kind must be set")
	}

	c.Author = strings.TrimSpace(c.Author)
	if c.Author == "" {
		return errors.ErrInvalid.WithMsgf("author must be set")
	}

	return c.genID()
}

func (c *Chunk) Apply(upd Updates) {
	if upd.Tags != nil {
		c.Tags = cleanTags(upd.Tags)
	}
}

func (c *Chunk) genID() error {
	if c.Kind == "" {
		return errors.ErrInternal.WithCausef("chunk kind not set")
	}

	idPrefix := strings.ToLower(string(c.Kind[0]))
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
