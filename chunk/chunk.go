package chunk

import (
	"regexp"
	"strings"
	"time"

	"github.com/chunked-app/cortex/pkg/errors"
)

const maxIDLen = 12

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// Chunk represents a piece of information written down by a user.
type Chunk struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Data      Data      `json:"data"`
	Tags      []string  `json:"tags"`
	Author    string    `json:"author"`
	Parent    string    `json:"parent,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	c.Type = c.Data.Type()

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
