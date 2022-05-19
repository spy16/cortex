package main

import (
	"encoding/json"
	"os"
)

func main() {
	m := map[string]interface{}{
		"foo": "bar",
	}

	b, _ := json.Marshal(m)

	c := chunk{Val: b}

	_ = json.NewEncoder(os.Stdout).Encode(c)
}

type chunk struct {
	Val json.RawMessage `json:"val"`
}
