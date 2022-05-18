package chunk

// Supported entity types within a chunk.
const (
	EntityTag     = "tag"
	EntityMention = "mention"
)

// Entity represents parsed information inside a chunk. This includes tags,
// mentions, links, etc.
type Entity struct {
	Type  string `json:"type"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

func parseEntities(str string) []Entity {
	var res []Entity
	// TODO: parse entities from the content of the chunk.
	return res
}
