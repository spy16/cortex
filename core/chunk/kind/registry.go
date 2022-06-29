package kind

import (
	"encoding/json"
	"strings"

	"github.com/chunked-app/cortex/pkg/errors"
)

type Registry struct{}

type chunkData interface {
	ValidateData() error
}

func (r *Registry) Validate(kind, data string) error {
	kind = strings.TrimSpace(strings.ToUpper(kind))

	var into chunkData
	switch kind {
	case "NOTE":
		into = &NoteData{}

	case "TODO":
		into = &TodoData{}

	case "IMAGE":
		into = &ImageData{}

	default:
		return errors.ErrInvalid.WithMsgf("invalid kind '%s'", kind)
	}

	if err := json.Unmarshal([]byte(data), into); err != nil {
		return errors.ErrInvalid.
			WithMsgf("failed to interpret data as '%s'", kind).
			WithCausef(err.Error())
	} else if err := into.ValidateData(); err != nil {
		return err
	}

	return nil
}
