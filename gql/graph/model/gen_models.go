// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// Chunk represents a piece of information.
type Chunk struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Rank      string    `json:"rank"`
	Data      string    `json:"data"`
	Tags      []string  `json:"tags"`
	AuthorID  string    `json:"author_id"`
	ParentID  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"Updated_at"`
}

// CreateRequest can be passed to createChunk mutation to create a new chunk.
type CreateRequest struct {
	ID       string  `json:"id"`
	Kind     Kind    `json:"kind"`
	Data     string  `json:"data"`
	Rank     string  `json:"rank"`
	AuthorID string  `json:"author_id"`
	ParentID *string `json:"parent_id"`
}

// UpdateRequest can be passed to updateChunk mutation to modify a chunk.
type UpdateRequest struct {
	Kind   *string `json:"kind"`
	Data   *string `json:"data"`
	Rank   *string `json:"rank"`
	Parent *string `json:"parent"`
}

// Kind represents the type of data in a chunk.
type Kind string

const (
	KindNote  Kind = "NOTE"
	KindImage Kind = "IMAGE"
	KindTodo  Kind = "TODO"
)

var AllKind = []Kind{
	KindNote,
	KindImage,
	KindTodo,
}

func (e Kind) IsValid() bool {
	switch e {
	case KindNote, KindImage, KindTodo:
		return true
	}
	return false
}

func (e Kind) String() string {
	return string(e)
}

func (e *Kind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Kind(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Kind", str)
	}
	return nil
}

func (e Kind) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}