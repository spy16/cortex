package chunk

import (
	"regexp"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

// Supported content-types in a chunk.
const (
	TypeText  = "text"
	TypeImage = "image"
	TypeLink  = "link"
)

const maxIDLen = 12

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// Chunk represents a piece of information written down by a user.
type Chunk struct {
	// chunk metadata.
	ID        string    `json:"id"`
	Author    string    `json:"author,omitempty"`
	Entities  []Entity  `json:"entities"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// content of the chunk.
	Type    string `json:"type"`
	Content string `json:"content"`

	// relations with other chunks.
	Parent      string `json:"parent,omitempty"`
	NextSibling string `json:"next_sibling,omitempty"`
	PrevSibling string `json:"prev_sibling,omitempty"`
}

func (c *Chunk) Validate() error {
	c.ID = strings.TrimSpace(c.ID)
	if err := validateID(c.ID, false); err != nil {
		return err
	}

	c.Parent = strings.TrimSpace(c.Parent)
	if err := validateID(c.Parent, true); err != nil {
		return err
	}

	c.Author = strings.TrimSpace(c.Author)
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
		c.UpdatedAt = c.CreatedAt
	}

	c.Content = strings.TrimSpace(c.Content)
	if c.Content == "" {
		return errors.ErrInvalid.WithMsgf("chunk data must not be empty")
	}

	c.Entities = parseEntities(c.Content)
	return nil
}

type Updates struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	ParentID    string `json:"parent_id"`
	PrevSibling string `json:"prev_sibling"`
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
