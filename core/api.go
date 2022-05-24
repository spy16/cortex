package core

import (
	"github.com/chunked-app/cortex/core/chunk"
	"github.com/chunked-app/cortex/core/user"
)

type API struct {
	Users  user.Store
	Chunks chunk.Store
}
